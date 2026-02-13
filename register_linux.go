//go:build linux

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

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	desktopDir := filepath.Join(home, ".local", "share", "applications")
	if err := os.MkdirAll(desktopDir, 0o755); err != nil {
		return err
	}

	desktop := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=Ope
Exec=%s %%u
Icon=ope
Terminal=false
Categories=Utility;
MimeType=x-scheme-handler/ope;
NoDisplay=true
`, exe)

	desktopFile := filepath.Join(desktopDir, "ope.desktop")
	if err := os.WriteFile(desktopFile, []byte(desktop), 0o644); err != nil {
		return err
	}

	// Register as default handler for ope:// scheme
	_ = exec.Command("xdg-mime", "default", "ope.desktop", "x-scheme-handler/ope").Run()
	_ = exec.Command("update-desktop-database", desktopDir).Run()

	fmt.Printf("Installed: %s\n", desktopFile)
	fmt.Println("URL scheme ope:// registered.")
	return nil
}

func uninstall() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	desktopFile := filepath.Join(home, ".local", "share", "applications", "ope.desktop")
	if err := os.Remove(desktopFile); err != nil && !os.IsNotExist(err) {
		return err
	}
	fmt.Printf("Removed: %s\n", desktopFile)
	return nil
}
