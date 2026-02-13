//go:build darwin

package main

import "os/exec"

func openPath(path string) error {
	return exec.Command("open", path).Run()
}
