package datatransfer

import (
	"encoding/json"
	"io"
	"net/http"
	"vortex/telehandler"
)

type Gofile struct {
	Status string `json:"status"`
	Data   struct {
		Server       string `json:"server"`
		Downloadpage string `json:"downloadPage"`
	} `json:"data"`
}

func gofileAvailableServer() string {
	var gofile Gofile

	req, _ := http.Get("https://api.gofile.io/getServer")
	req.Header.Set("User-Agent", telehandler.UserAgent)
	msgBody, _ := io.ReadAll(req.Body)
	json.Unmarshal(msgBody, &gofile)

	return gofile.Data.Server
}
