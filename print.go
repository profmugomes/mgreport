// Copyright (C) 2026 Murilo Gomes <profmugomes.com.br>
// SPDX-License-Identifier: MIT

// Site: https://www.profmugomes.com.br

package mgreport

import (
	"os"
	"os/exec"
	"runtime"
	"time"
)

// PrintReport generates a temporary PDF and opens it with the system's
// default PDF viewer, allowing the user to print from there.
func PrintReport(report *Report) error {
	tmpDir := os.TempDir()
	tmpFile := tmpDir + "/mgreport_print_" + time.Now().Format("20060102150405") + ".pdf"

	err := ExportToPDF(report, tmpFile)
	if err != nil {
		return err
	}

	return openFile(tmpFile)
}

// openFile opens a file with the default system application.
func openFile(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", path)
	case "darwin":
		cmd = exec.Command("open", path)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", path)
	default:
		return os.ErrInvalid
	}

	return cmd.Start()
}
