package model

import (
	"encoding/json"
	"os/exec"
	"strings"
)

type Device struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func ParseDevices(bytes []byte) ([]Device, error) {
	var devices []Device
	err := json.Unmarshal(bytes, &devices)

	if err != nil {
		return devices, err
	}

	return devices, nil
}

func ParseEmulators(bytes []byte) ([]Device, error) {
	var devices []Device

	lines := strings.Split(string(bytes), "\n")

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

// Whether device has sufficent data
func (d Device) Verified() bool {
	return d.Name != "" && d.ID != ""
}

func (d Device) BuildLaunchEmulatorCommand(isCold bool) *exec.Cmd {
	args := []string{"emulators", "--launch", d.Name}

	if isCold {
		args = append(args, "--cold")
	}

	c := exec.Command("flutter", args...)

	return c
}

func FlutterEmulators() *exec.Cmd {
	return exec.Command("flutter", "emulators")
}

func FlutterDevices() *exec.Cmd {
	return exec.Command("flutter", "devices", "--machine")
}
