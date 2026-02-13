//go:build darwin

package main

import (
	"os/exec"
	"strings"
)

// ConfirmResult represents the user's choice in a confirmation dialog.
type ConfirmResult int

const (
	ConfirmAllow  ConfirmResult = iota
	ConfirmAlways
	ConfirmBlock
	ConfirmCancel
)

func showErrorDialog(title, message string) {
	script := `display dialog "` + escapeAS(message) + `" with title "` + escapeAS(title) + `" buttons {"OK"} default button "OK" with icon stop`
	_ = exec.Command("osascript", "-e", script).Run()
}

func showConfirmDialog(path string) ConfirmResult {
	script := `display dialog "Open this path?\n\n` + escapeAS(path) + `" ` +
		`with title "ope — Confirm" ` +
		`buttons {"Block", "Always Allow", "Allow Once"} ` +
		`default button "Allow Once" ` +
		`with icon caution`
	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return ConfirmCancel
	}
	result := strings.TrimSpace(string(out))
	switch {
	case strings.Contains(result, "Allow Once"):
		return ConfirmAllow
	case strings.Contains(result, "Always Allow"):
		return ConfirmAlways
	case strings.Contains(result, "Block"):
		return ConfirmBlock
	default:
		return ConfirmCancel
	}
}

func escapeAS(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}
