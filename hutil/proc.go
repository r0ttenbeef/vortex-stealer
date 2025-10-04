//go:build windows

package hutil

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"github.com/yusufpapurcu/wmi"
	"golang.org/x/sys/windows"
)

const (
	TH32CS_SNAPPROCESS        = 0x00000002
	PROCESS_QUERY_INFORMATION = 0x0400
)

type WindowsProcess struct {
	ProcessID         int
	ParentProcessID   int
	ExecutableProcess string
}

type Win32_Process struct {
	Name           string
	ExecutablePath string
}

func getProcessList() ([]WindowsProcess, error) {
	var procEntry windows.ProcessEntry32
	procList := make([]WindowsProcess, 0, 50)

	pHandle, err := windows.CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(pHandle)

	procEntry.Size = uint32(unsafe.Sizeof(procEntry))

	if err = windows.Process32First(pHandle, &procEntry); err != nil {
		return nil, err
	}

	for {
		procList = append(procList, newWindowsProcss(&procEntry))
		if err = windows.Process32Next(pHandle, &procEntry); err != nil {
			if err == syscall.ERROR_NO_MORE_FILES {
				return procList, nil
			}
			return nil, err
		}
	}
}

func newWindowsProcss(e *windows.ProcessEntry32) WindowsProcess {
	var end = 0

	for {
		if e.ExeFile[end] == 0 {
			break
		}
		end++
	}

	return WindowsProcess{
		ProcessID:         int(e.ProcessID),
		ParentProcessID:   int(e.ParentProcessID),
		ExecutableProcess: syscall.UTF16ToString(e.ExeFile[:end]),
	}
}

// Get process ID from process name
func GetProcessIdByName(procName string) (int, error) {
	procs, err := getProcessList()
	if err != nil {
		return 0, err
	}

	for _, p := range procs {
		if strings.Contains(strings.ToLower(p.ExecutableProcess), strings.ToLower(procName)) {
			return p.ProcessID, nil
		}
	}

	return 0, nil
}

// Get Running process path location
func GetProcessPath(pName string) (string, error) {
	var path []Win32_Process

	query := fmt.Sprint("WHERE Name = '%s'", pName)
	q := wmi.CreateQuery(&path, query)

	if err := wmi.Query(q, &path); err != nil {
		return "", err
	}

	for _, v := range path {
		if v.Name == pName {
			return v.ExecutablePath, nil
		}
	}

	return "", nil
}

// Terminate running process by name
func TerminateProcess(procName string) error {
	procID, err := GetProcessIdByName(procName)
	if err != nil {
		return err
	}

	if procID != 0 {
		pHandle, _ := syscall.OpenProcess(syscall.PROCESS_TERMINATE, false, uint32(procID))
		if err = syscall.TerminateProcess(pHandle, 1); err != nil {
			return err
		}
	}

	return nil
}
