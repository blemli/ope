//go:build windows

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func install() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot find executable: %w", err)
	}
	exe, _ = filepath.Abs(exe)

	// Create a VBS launcher script next to the exe that runs it without a console window.
	// WshShell.Run with 0 = hidden window.
	launcherPath := filepath.Join(filepath.Dir(exe), "ope-launcher.vbs")
	vbs := fmt.Sprintf(`Set WshShell = CreateObject("WScript.Shell")`+"\r\n"+
		`WshShell.Run Chr(34) & "%s" & Chr(34) & " " & Chr(34) & WScript.Arguments(0) & Chr(34), 0, False`+"\r\n",
		exe)
	if err := os.WriteFile(launcherPath, []byte(vbs), 0o644); err != nil {
		return fmt.Errorf("cannot write launcher: %w", err)
	}

	// Register ope:// URL scheme in HKCU (no admin required)
	// The handler invokes the VBS launcher to avoid a console window flash
	handler := fmt.Sprintf(`wscript.exe "%s" "%%1"`, launcherPath)
	cmds := [][]string{
		{"reg", "add", `HKCU\Software\Classes\ope`, "/ve", "/d", "URL:ope Protocol", "/f"},
		{"reg", "add", `HKCU\Software\Classes\ope`, "/v", "URL Protocol", "/d", "", "/f"},
		{"reg", "add", `HKCU\Software\Classes\ope\shell\open\command`, "/ve", "/d", handler, "/f"},
	}

	for _, args := range cmds {
		if err := exec.Command(args[0], args[1:]...).Run(); err != nil {
			return fmt.Errorf("registry command failed: %w", err)
		}
	}

	fmt.Println("URL scheme ope:// registered in Windows registry.")
	return nil
}

func uninstall() error {
	// Remove VBS launcher
	exe, _ := os.Executable()
	exe, _ = filepath.Abs(exe)
	launcherPath := filepath.Join(filepath.Dir(exe), "ope-launcher.vbs")
	_ = os.Remove(launcherPath)

	err := exec.Command("reg", "delete", `HKCU\Software\Classes\ope`, "/f").Run()
	if err != nil {
		return fmt.Errorf("failed to remove registry keys: %w", err)
	}
	fmt.Println("URL scheme ope:// unregistered.")
	return nil
}
