package ui

import (
	"os/exec"
	"runtime"
	"strings"
)

// readClipboard returns the current clipboard text, or "" on failure.
func readClipboard() string {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("powershell", "-NoProfile", "-Command", "Get-Clipboard")
	case "darwin":
		cmd = exec.Command("pbpaste")
	default:
		// Try xclip, then xsel.
		if p, err := exec.LookPath("xclip"); err == nil {
			cmd = exec.Command(p, "-selection", "clipboard", "-o")
		} else if p, err := exec.LookPath("xsel"); err == nil {
			cmd = exec.Command(p, "--clipboard", "--output")
		} else {
			return ""
		}
	}
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
