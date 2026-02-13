//go:build windows

package main

import (
	"os"
	"os/exec"
)

func openPath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return exec.Command("explorer", path).Run()
	}
	return exec.Command("cmd", "/c", "start", "", path).Run()
}
