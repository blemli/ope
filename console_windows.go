//go:build windows

package main

import (
	"os"
	"syscall"
	"unsafe"
)

// attachConsole attaches to the parent process's console so CLI output works
// even though the binary is built with -H windowsgui (GUI subsystem).
// When launched from a browser (no parent console), this is a no-op.
func init() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	attachConsole := kernel32.NewProc("AttachConsole")
	getConsoleMode := kernel32.NewProc("GetConsoleMode")

	const ATTACH_PARENT_PROCESS = ^uintptr(0) // -1

	r, _, _ := attachConsole.Call(ATTACH_PARENT_PROCESS)
	if r == 0 {
		return // no parent console (launched from browser) — nothing to do
	}

	// Reopen stdout/stderr to the attached console
	getStdHandle := kernel32.NewProc("GetStdHandle")

	const STD_OUTPUT_HANDLE = ^uintptr(0) - 10 + 1 // -11
	const STD_ERROR_HANDLE = ^uintptr(0) - 10       // -12

	hOut, _, _ := getStdHandle.Call(STD_OUTPUT_HANDLE)
	hErr, _, _ := getStdHandle.Call(STD_ERROR_HANDLE)

	// Check if stdout is a valid console handle
	var mode uint32
	r, _, _ = getConsoleMode.Call(hOut, uintptr(unsafe.Pointer(&mode)))
	if r != 0 {
		os.Stdout = os.NewFile(hOut, "stdout")
		os.Stderr = os.NewFile(hErr, "stderr")
	}
}
