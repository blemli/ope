# Plan: `ope` — Cross-platform URL scheme handler

Rewrite OpenFolderApp in Go as `ope`, supporting macOS, Windows, and Linux.
URL scheme: `ope://` (replaces `openfolder://`)

## Project Structure

```
ope/
├── main.go                  # Entry point: parse URL or run subcommands
├── config.go                # Load/save ope.yml, blocklist/allowlist logic
├── open.go                  # Shared: path expansion, glob, validation
├── open_darwin.go           # macOS: `open <path>`
├── open_windows.go          # Windows: `explorer` / `cmd /c start`
├── open_linux.go            # Linux: `xdg-open <path>`
├── dialog_darwin.go         # macOS: osascript alert
├── dialog_windows.go        # Windows: PowerShell MessageBox
├── dialog_linux.go          # Linux: zenity / notify-send fallback
├── register_darwin.go       # macOS: no-op (Info.plist handles it)
├── register_windows.go      # Windows: add/remove registry keys
├── register_linux.go        # Linux: install .desktop file + xdg-mime
├── packaging/
│   ├── macos/
│   │   ├── Info.plist       # URL scheme + app metadata
│   │   └── entitlements.plist
│   ├── linux/
│   │   └── ope.desktop      # XDG desktop entry
│   └── icon.icns / icon.ico / icon.png  # User-provided icons
├── Makefile                 # Build + package for all platforms
├── .goreleaser.yml          # Automated cross-platform releases
├── .github/workflows/
│   └── release.yml          # CI: build on tag push, create GitHub release
├── go.mod
└── README.md
```

## Implementation Steps

### Step 1: Initialize Go project
- `go mod init github.com/blemli/ope`
- Basic main.go with arg parsing

### Step 2: Core logic (open.go)
Shared across all platforms:
- Parse `ope://` URL → extract path (strip scheme, handle `ope:///absolute` and `ope://~/relative`)
- Tilde expansion (`~` → home dir)
- Wildcard/glob expansion (`*` patterns)
- Path validation (exists? file or dir?)
- **Security check**: validate path/extension against config before opening

### Step 2b: Config & security (config.go)
Single config file: `ope.yml` (located in OS config dir, e.g. `~/.config/ope/ope.yml`)

```yaml
# ope.yml
blocked:          # Always blocked — never opened, no prompt
  extensions:
    - .exe
    - .bat
    - .cmd
    - .ps1
    - .sh
    - .msi
    - .scr
    - .vbs
    - .wsf
  paths:
    - /Windows/System32
    - /etc

allowed:          # Previously approved — opened without prompt
  extensions:
    - .pdf
    - .txt
    - .docx
    - .xlsx
    - .png
    - .jpg
  paths:
    - ~/Documents
    - ~/Downloads
```

**Security flow when opening a path:**
1. Check `blocked` list → if matched, show error dialog "Blocked by policy", refuse
2. Check `allowed` list → if matched, open immediately
3. Otherwise → show confirmation dialog: "Open <path>? [Allow once / Always allow / Block]"
   - "Allow once": open this time only
   - "Always allow": add extension/path to `allowed` in ope.yml, then open
   - "Block": add to `blocked` in ope.yml, refuse

### Step 3: Platform openers (open_*.go, build-tagged)
- **macOS**: `exec.Command("open", path)` — works for both files and dirs
- **Windows**: `exec.Command("explorer", path)` for dirs, `exec.Command("cmd", "/c", "start", "", path)` for files
- **Linux**: `exec.Command("xdg-open", path)`

### Step 4: Error dialogs (dialog_*.go, build-tagged)
Show native error when path doesn't exist:
- **macOS**: `osascript -e 'display dialog "..."'`
- **Windows**: PowerShell `[System.Windows.MessageBox]::Show(...)`
- **Linux**: `zenity --error` with fallback to `notify-send`

### Step 5: CLI subcommands
`ope` binary supports:
- `ope <url>` — default: handle the URL (called by OS)
- `ope install` — register URL scheme for current platform
- `ope uninstall` — remove URL scheme registration
- `ope config` — print config file path and current settings
- `ope version` — print version

**Registration**: runs automatically as post-install hook in package managers:
- **Homebrew**: `postflight` runs `ope install`
- **Scoop**: `post_install` runs `ope install`
- **Linux install.sh**: runs `ope install` at end of script
- **Fallback**: if invoked without registration, prompt: "ope:// protocol is not registered. Register now? [y/n]"

### Step 6: URL scheme registration (register_*.go)

**macOS** (`ope install`):
- Create .app bundle at `~/Applications/Ope.app/`
- Copy binary into `Contents/MacOS/ope`
- Write `Info.plist` with CFBundleURLTypes for `ope://`
- Register with Launch Services: `/System/Library/Frameworks/CoreServices.framework/Versions/A/Frameworks/LaunchServices.framework/Versions/A/Support/lsregister -R ~/Applications/Ope.app`

**Windows** (`ope install`):
- Add registry keys:
  ```
  HKCU\Software\Classes\ope
    (Default) = "URL:ope Protocol"
    "URL Protocol" = ""
    shell\open\command\(Default) = "C:\path\to\ope.exe" "%1"
  ```
- Uses HKCU (no admin needed)

**Linux** (`ope install`):
- Write `ope.desktop` to `~/.local/share/applications/`
- Run `xdg-mime default ope.desktop x-scheme-handler/ope`
- Run `update-desktop-database ~/.local/share/applications/`

### Step 7: macOS .app bundle
For Homebrew distribution, we ship a .app bundle (zip):
```
Ope.app/
  Contents/
    Info.plist      # CFBundleURLTypes, LSBackgroundOnly
    MacOS/
      ope           # Go binary
    Resources/
      icon.icns     # App icon
```
The Makefile creates this structure and zips it.

### Step 8: Build system (Makefile)
```makefile
build-macos:    # go build → create .app bundle → zip
build-windows:  # GOOS=windows go build → zip with .exe
build-linux:    # GOOS=linux go build → tarball with binary + .desktop + install.sh
```

### Step 9: Goreleaser + GitHub Actions
- `.goreleaser.yml`: cross-compile, create archives, homebrew tap
- On git tag push → build all platforms → create GitHub release with assets
- Post-build hooks to create .app bundle for macOS

### Step 10: Package manifests

**Homebrew Cask** (`homebrew-tap/Casks/ope.rb`):
- Downloads macOS .app zip from GitHub releases
- Moves Ope.app to /Applications
- Runs `ope install` (lsregister) as postflight

**Scoop** (`blemli-bucket/bucket/ope.json`):
- Downloads Windows .exe zip from GitHub releases
- Runs `ope install` as post_install (registry)
- Runs `ope uninstall` as pre_uninstall

**Linux**:
- Homebrew tap works on Linuxbrew too (formula, not cask)
- Also provide tarball with install.sh script
- install.sh: copies binary to /usr/local/bin, runs `ope install`

### Step 11: Testing
- Update test.html with `ope://` links
- Manual testing on each platform
- Go unit tests for URL parsing and path expansion

## Summary of user-facing commands
```bash
# Install
brew install blemli/tap/ope          # macOS (also Linuxbrew)
scoop bucket add blemli https://github.com/blemli/blemli-bucket
scoop install ope                     # Windows
# Linux manual: download tarball, run install.sh

# Usage — click ope:// links in browser, they just work
# Or from terminal:
ope "ope:///Users/you/Documents"
ope install     # register URL scheme
ope uninstall   # unregister
```

## Decisions made
- [x] **Repo**: New repo at `github.com/blemli/ope` (user creates it)
- [x] **First run**: Ask user to register URL scheme on first run (not silent)
- [x] **Config**: Single `ope.yml` file, no GUI
- [x] **Security**: Blocklist for dangerous types, confirmation dialog for unknown, auto-add to allowed/blocked based on user choice

## Open items
- [ ] User to provide icon (will need .icns for macOS, .ico for Windows, .png for Linux)
