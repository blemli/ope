# ope

Open files and folders from the browser via `ope://` URLs.

For example [ope://~/Downloads](ope://~/Downloads)

<img src="ope.svg" width="128" />

## Install

### macOS (Homebrew)

```bash
brew install blemli/tap/ope
```

### Windows (Scoop)

```powershell
scoop bucket add blemli https://github.com/blemli/blemli-bucket
scoop install ope
```

### Linux

Download the binary from [Releases](https://github.com/blemli/ope/releases) and run:

```bash
chmod +x ope
sudo cp ope /usr/local/bin/
ope install
```

### From source

```bash
go install github.com/blemli/ope@latest
ope install
```

## Usage

Open a folder in your file manager:

```
ope:///Users/me/Documents
```

Open a specific file:

```
ope:///Users/me/notes.txt
```

Use in HTML links:

```html
<a href="ope:///tmp">Open /tmp</a>
```

## CLI

```
ope <ope://url>    Open a file or folder
ope install        Register ope:// URL scheme
ope uninstall      Unregister ope:// URL scheme
ope config         Show configuration
ope version        Print version
```

## Security

On first use, `ope` creates a config file with a blocklist of dangerous extensions (`.exe`, `.bat`, `.cmd`, etc.). When opening an unknown file type, a confirmation dialog asks you to:

- **Allow Once** — open this time only
- **Always Allow** — add to allowlist
- **Block** — add to blocklist

Config location:
- macOS: `~/Library/Application Support/ope/ope.yml`
- Windows: `%APPDATA%\ope\ope.yml`
- Linux: `~/.config/ope/ope.yml`

## Building

```bash
make build          # Build for current platform
make build-macos    # Universal macOS binary + Ope.app
make build-windows  # Windows binary
make build-linux    # Linux binaries (amd64 + arm64)
make test           # Run tests
make icons          # Generate icons from SVG
```

## License

MIT
