package helpers

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func GetSSID() (string, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("netsh", "wlan", "show", "interfaces")
	case "darwin":
		cmd = exec.Command("/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport", "-I")
	case "linux":
		cmd = exec.Command("iwgetid", "-r")
	default:
		return "", fmt.Errorf("unsupported OS")
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	output := out.String()
	switch runtime.GOOS {
	case "windows":
		for _, line := range strings.Split(output, "\n") {
			if strings.Contains(line, "SSID") && !strings.Contains(line, "BSSID") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					return strings.TrimSpace(parts[1]), nil
				}
			}
		}
	case "darwin":
		for _, line := range strings.Split(output, "\n") {
			if strings.Contains(line, " SSID:") {
				parts := strings.Split(line, ": ")
				if len(parts) > 1 {
					return strings.TrimSpace(parts[1]), nil
				}
			}
		}
	case "linux":
		return strings.TrimSpace(output), nil
	}

	return "", fmt.Errorf("SSID not found")
}
