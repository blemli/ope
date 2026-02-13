//go:build linux

package main

import "os/exec"

func openPath(path string) error {
	return exec.Command("xdg-open", path).Run()
}
