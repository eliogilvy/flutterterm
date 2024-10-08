package utils

import (
	"encoding/json"
	"os/exec"
	"strings"
)

type Device struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

// Returns a list of available devices to use
func GetDevices() ([]Device, error) {
	var devices []Device

	cmd := exec.Command("flutter", "devices", "--machine")
	output, err := cmd.Output()
	if err != nil {
		PrintError(err.Error())
		return devices, err
	}

	err = json.Unmarshal(output, &devices)

	if err != nil {
		return devices, err
	}

	return devices, nil
}

// Return a list of available flutter emulators
func GetEmulators() ([]Device, error) {
	PrintInfo("Getting emulators\n\n")
	var devices []Device

	cmd := exec.Command("flutter", "emulators")

	output, err := cmd.Output()

	if err != nil {
		return devices, err
	}

	lines := strings.Split(string(output), "\n")

	for i, line := range lines {
		if line == "" {
			continue
		}
		// No useful info on these lines
		if i >= 0 && i < 3 {
			continue
		}

		// Emulators start on line 4

		if line == "" {
			break
		}

		parts := strings.Split(line, "•")

		if len(parts) < 4 {
			continue
		}

		device := Device{
			ID:   strings.TrimSpace(parts[0]),
			Name: strings.TrimSpace(parts[1]),
		}

		devices = append(devices, device)
	}

	// Remove the first element which is "Name"
	devices = devices[0:]

	return devices, nil
}
