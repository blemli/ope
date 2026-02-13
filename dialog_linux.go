//go:build linux

package main

import (
	"os/exec"
	"strings"
)

type ConfirmResult int

const (
	ConfirmAllow  ConfirmResult = iota
	ConfirmAlways
	ConfirmBlock
	ConfirmCancel
)

func showErrorDialog(title, message string) {
	// Try zenity first, fall back to notify-send
	err := exec.Command("zenity", "--error", "--title="+title, "--text="+message).Run()
	if err != nil {
		_ = exec.Command("notify-send", "-u", "critical", title, message).Run()
	}
}

func showConfirmDialog(path string) ConfirmResult {
	// Use zenity --list for a 3-option dialog
	out, err := exec.Command("zenity", "--list",
		"--title=ope — Confirm",
		"--text=Open this path?\n\n"+path,
		"--column=Action",
		"Allow Once",
		"Always Allow",
		"Block",
	).Output()
	if err != nil {
		return ConfirmCancel
	}
	result := strings.TrimSpace(string(out))
	switch result {
	case "Allow Once":
		return ConfirmAllow
	case "Always Allow":
		return ConfirmAlways
	case "Block":
		return ConfirmBlock
	default:
		return ConfirmCancel
	}
}
