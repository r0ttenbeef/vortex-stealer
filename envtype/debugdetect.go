//go:build windows

package envtype

import (
	"syscall"
	"unsafe"
)

// Check if running inside debugger
func DetectDebugging() (bool, error) {
	var isDebugger bool

	kernel32, err := syscall.LoadLibrary("kernel32.dll")
	if err != nil {
		return false, err
	}
	defer syscall.FreeLibrary(kernel32)

	isDebuggerPresent, err := syscall.GetProcAddress(kernel32, "IsDebuggerPresent")
	if err != nil {
		return false, err
	}

	syscall.SyscallN(isDebuggerPresent, 0, uintptr(unsafe.Pointer(&isDebugger)), 0, 0)

	if isDebugger {
		return true, nil
	} else {
		return false, nil
	}
}
