//go:build windows

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
	ps := `Add-Type -AssemblyName System.Windows.Forms; ` +
		`[System.Windows.Forms.MessageBox]::Show("` + escapePSStr(message) + `", "` + escapePSStr(title) + `", "OK", "Error")`
	_ = exec.Command("powershell", "-NoProfile", "-Command", ps).Run()
}

func showConfirmDialog(path string) ConfirmResult {
	// PowerShell script that shows a custom form with 3 buttons
	ps := `Add-Type -AssemblyName System.Windows.Forms
$form = New-Object System.Windows.Forms.Form
$form.Text = "ope - Confirm"
$form.Width = 450
$form.Height = 200
$form.StartPosition = "CenterScreen"
$form.FormBorderStyle = "FixedDialog"
$form.MaximizeBox = $false

$label = New-Object System.Windows.Forms.Label
$label.Text = "Open this path?` + "\n\n" + escapePSStr(path) + `"
$label.AutoSize = $true
$label.Location = New-Object System.Drawing.Point(20, 20)
$form.Controls.Add($label)

$btnAllow = New-Object System.Windows.Forms.Button
$btnAllow.Text = "Allow Once"
$btnAllow.Location = New-Object System.Drawing.Point(100, 120)
$btnAllow.Add_Click({ $form.Tag = "allow"; $form.Close() })
$form.Controls.Add($btnAllow)

$btnAlways = New-Object System.Windows.Forms.Button
$btnAlways.Text = "Always Allow"
$btnAlways.Location = New-Object System.Drawing.Point(200, 120)
$btnAlways.Add_Click({ $form.Tag = "always"; $form.Close() })
$form.Controls.Add($btnAlways)

$btnBlock = New-Object System.Windows.Forms.Button
$btnBlock.Text = "Block"
$btnBlock.Location = New-Object System.Drawing.Point(310, 120)
$btnBlock.Add_Click({ $form.Tag = "block"; $form.Close() })
$form.Controls.Add($btnBlock)

$form.ShowDialog() | Out-Null
$form.Tag`

	out, err := exec.Command("powershell", "-NoProfile", "-Command", ps).Output()
	if err != nil {
		return ConfirmCancel
	}
	result := strings.TrimSpace(string(out))
	switch result {
	case "allow":
		return ConfirmAllow
	case "always":
		return ConfirmAlways
	case "block":
		return ConfirmBlock
	default:
		return ConfirmCancel
	}
}

func escapePSStr(s string) string {
	s = strings.ReplaceAll(s, "`", "``")
	s = strings.ReplaceAll(s, `"`, "`\"")
	s = strings.ReplaceAll(s, `$`, "`$")
	return s
}
