package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ParseOpeURL extracts the local file path from an ope:// URL.
// Supports: ope:///path, ope://localhost/path, ope://path
func ParseOpeURL(raw string) (string, error) {
	// Handle ope:path (no slashes) as ope:///path
	if strings.HasPrefix(raw, "ope:") && !strings.HasPrefix(raw, "ope://") {
		raw = "ope:///" + strings.TrimPrefix(raw, "ope:")
	}

	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	if u.Scheme != "ope" {
		return "", fmt.Errorf("unsupported scheme: %s", u.Scheme)
	}

	path := u.Host + u.Path
	if u.Host == "localhost" {
		path = u.Path
	}

	// URL-decode the path
	path, err = url.PathUnescape(path)
	if err != nil {
		return "", fmt.Errorf("invalid path encoding: %w", err)
	}

	// Windows drive letter: /C:/... → C:/...
	if runtime.GOOS == "windows" && len(path) >= 3 && path[0] == '/' && path[2] == ':' {
		path = path[1:]
	}

	return path, nil
}

// ExpandPath handles tilde expansion and glob patterns.
func ExpandPath(path string) (string, error) {
	// Tilde expansion
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot expand ~: %w", err)
		}
		path = filepath.Join(home, path[1:])
	}

	// Glob expansion — if the path contains wildcards, resolve them
	if strings.ContainsAny(path, "*?[") {
		matches, err := filepath.Glob(path)
		if err != nil {
			return "", fmt.Errorf("invalid glob pattern: %w", err)
		}
		if len(matches) == 0 {
			return "", fmt.Errorf("no files matched: %s", path)
		}
		// Use the first match
		path = matches[0]
	}

	return filepath.Clean(path), nil
}

// HandleURL is the main entry point: parse URL, expand path, check security, open.
func HandleURL(raw string) error {
	path, err := ParseOpeURL(raw)
	if err != nil {
		showErrorDialog("Invalid URL", err.Error())
		return err
	}

	path, err = ExpandPath(path)
	if err != nil {
		showErrorDialog("Path Error", err.Error())
		return err
	}

	// Check that path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		msg := fmt.Sprintf("Path does not exist: %s", path)
		showErrorDialog("Not Found", msg)
		return fmt.Errorf("path does not exist: %s", path)
	}

	cfg, err := LoadConfig()
	if err != nil {
		showErrorDialog("Config Error", err.Error())
		return err
	}

	switch cfg.CheckSecurity(path) {
	case ActionBlock:
		msg := fmt.Sprintf("Blocked by security policy: %s", filepath.Base(path))
		showErrorDialog("Blocked", msg)
		return fmt.Errorf("blocked by security policy: %s", filepath.Base(path))

	case ActionAllow:
		return openPath(path)

	case ActionAsk:
		result := showConfirmDialog(path)
		switch result {
		case ConfirmAllow:
			return openPath(path)
		case ConfirmAlways:
			cfg.Allowed = append(cfg.Allowed, filepath.Base(path))
			_ = SaveConfig(cfg)
			return openPath(path)
		case ConfirmBlock:
			cfg.Blocked = append(cfg.Blocked, filepath.Base(path))
			_ = SaveConfig(cfg)
			return fmt.Errorf("blocked: %s", filepath.Base(path))
		default:
			return fmt.Errorf("cancelled")
		}
	}

	return nil
}
