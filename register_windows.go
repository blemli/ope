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

	// Register ope:// URL scheme in HKCU (no admin required)
	cmds := [][]string{
		{"reg", "add", `HKCU\Software\Classes\ope`, "/ve", "/d", "URL:ope Protocol", "/f"},
		{"reg", "add", `HKCU\Software\Classes\ope`, "/v", "URL Protocol", "/d", "", "/f"},
		{"reg", "add", `HKCU\Software\Classes\ope\shell\open\command`, "/ve", "/d", fmt.Sprintf(`"%s" "%%1"`, exe), "/f"},
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
	err := exec.Command("reg", "delete", `HKCU\Software\Classes\ope`, "/f").Run()
	if err != nil {
		return fmt.Errorf("failed to remove registry keys: %w", err)
	}
	fmt.Println("URL scheme ope:// unregistered.")
	return nil
}
