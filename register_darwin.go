//go:build darwin

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func install() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot find executable: %w", err)
	}
	exe, _ = filepath.EvalSymlinks(exe)

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	appDir := filepath.Join(home, "Applications", "Ope.app")
	contentsDir := filepath.Join(appDir, "Contents")
	macosDir := filepath.Join(contentsDir, "MacOS")
	resourcesDir := filepath.Join(contentsDir, "Resources")

	// Create directory structure
	for _, dir := range []string{macosDir, resourcesDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("cannot create %s: %w", dir, err)
		}
	}

	// Copy binary
	input, err := os.ReadFile(exe)
	if err != nil {
		return fmt.Errorf("cannot read binary: %w", err)
	}
	binDst := filepath.Join(macosDir, "ope")
	if err := os.WriteFile(binDst, input, 0o755); err != nil {
		return fmt.Errorf("cannot write binary: %w", err)
	}

	// Write Info.plist
	plist := strings.ReplaceAll(infoPlistTemplate, "VERSION", Version)
	if err := os.WriteFile(filepath.Join(contentsDir, "Info.plist"), []byte(plist), 0o644); err != nil {
		return fmt.Errorf("cannot write Info.plist: %w", err)
	}

	// Copy icon if available
	iconSrc := filepath.Join(filepath.Dir(exe), "..", "packaging", "icons", "icon.icns")
	if iconData, err := os.ReadFile(iconSrc); err == nil {
		_ = os.WriteFile(filepath.Join(resourcesDir, "icon.icns"), iconData, 0o644)
	}

	// Register URL scheme
	_ = exec.Command("/System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister",
		"-R", appDir).Run()

	fmt.Printf("Installed: %s\n", appDir)
	fmt.Println("URL scheme ope:// registered.")
	return nil
}

func uninstall() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	appDir := filepath.Join(home, "Applications", "Ope.app")
	if err := os.RemoveAll(appDir); err != nil {
		return err
	}
	fmt.Printf("Removed: %s\n", appDir)
	return nil
}

const infoPlistTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleIdentifier</key>
	<string>li.blem.ope</string>
	<key>CFBundleName</key>
	<string>Ope</string>
	<key>CFBundleDisplayName</key>
	<string>Ope</string>
	<key>CFBundleVersion</key>
	<string>VERSION</string>
	<key>CFBundleShortVersionString</key>
	<string>VERSION</string>
	<key>CFBundleExecutable</key>
	<string>ope</string>
	<key>CFBundleIconFile</key>
	<string>icon</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
	<key>LSBackgroundOnly</key>
	<true/>
	<key>CFBundleURLTypes</key>
	<array>
		<dict>
			<key>CFBundleURLName</key>
			<string>li.blem.ope</string>
			<key>CFBundleURLSchemes</key>
			<array>
				<string>ope</string>
			</array>
		</dict>
	</array>
</dict>
</plist>`
