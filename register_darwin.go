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

	// Remove old bundle
	_ = os.RemoveAll(appDir)

	// Create directory structure
	for _, dir := range []string{macosDir, resourcesDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("cannot create %s: %w", dir, err)
		}
	}

	// Compile the Swift launcher that handles Apple Events
	launcherSrc := filepath.Join(filepath.Dir(exe), "..", "packaging", "macos", "launcher.swift")
	launcherDst := filepath.Join(macosDir, "ope-launcher")

	// Try to find the launcher source; if not found, use embedded source
	if _, err := os.Stat(launcherSrc); os.IsNotExist(err) {
		// Write embedded launcher source to temp file
		tmpSrc := filepath.Join(os.TempDir(), "ope-launcher.swift")
		if err := os.WriteFile(tmpSrc, []byte(launcherSwiftSource), 0o644); err != nil {
			return fmt.Errorf("cannot write launcher source: %w", err)
		}
		launcherSrc = tmpSrc
		defer os.Remove(tmpSrc)
	}

	cmd := exec.Command("swiftc", "-O", "-o", launcherDst, launcherSrc)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("swiftc failed: %w\n%s", err, string(out))
	}

	// Copy our Go binary into Resources
	input, err := os.ReadFile(exe)
	if err != nil {
		return fmt.Errorf("cannot read binary: %w", err)
	}
	if err := os.WriteFile(filepath.Join(resourcesDir, "ope"), input, 0o755); err != nil {
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

	// Register URL scheme with LaunchServices
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
	<string>ope-launcher</string>
	<key>CFBundleIconFile</key>
	<string>icon</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
	<key>LSUIElement</key>
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

const launcherSwiftSource = `import Cocoa

class AppDelegate: NSObject, NSApplicationDelegate {
    func applicationDidFinishLaunching(_ notification: Notification) {
        NSAppleEventManager.shared().setEventHandler(
            self,
            andSelector: #selector(handleURL(_:withReply:)),
            forEventClass: AEEventClass(kInternetEventClass),
            andEventID: AEEventID(kAEGetURL)
        )
    }

    @objc func handleURL(_ event: NSAppleEventDescriptor, withReply reply: NSAppleEventDescriptor) {
        guard let urlString = event.paramDescriptor(forKeyword: keyDirectObject)?.stringValue else {
            return
        }

        let bundle = Bundle.main
        let binPath = bundle.resourceURL!.appendingPathComponent("ope").path

        let task = Process()
        task.executableURL = URL(fileURLWithPath: binPath)
        task.arguments = [urlString]
        try? task.run()
        task.waitUntilExit()

        DispatchQueue.main.asyncAfter(deadline: .now() + 0.5) {
            NSApplication.shared.terminate(nil)
        }
    }
}

let app = NSApplication.shared
let delegate = AppDelegate()
app.delegate = delegate
app.setActivationPolicy(.accessory)
app.run()
`
