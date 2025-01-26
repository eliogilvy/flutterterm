package utils

import (
	"encoding/json"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

type Device struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func (d Device) FilterValue() string {
	return d.Name
}

func (d Device) Title() string {
	return d.Name
}

func (d Device) Description() string {
	return d.ID
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

func GetDeviceList(devices []Device) list.Model {
	var items []list.Item
	for _, device := range devices {
		item := device
		items = append(items, item)
	}

	l := GetList(items)
	l.Title = "Devices"
	l.SetSize(100, 30)

	return l
}

func GetDeviceTable(devices []Device) table.Model {
	c := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "ID", Width: 20},
	}

	var r []table.Row

	return table.New(
		table.WithColumns(c),
		table.WithRows(r),
		table.WithFocused(true),
		// table.WithHeight(len(devices)+1),
		table.WithHeight(15),
		// Cool style
		table.WithStyles(table.Styles{
			// Cool header
			Header: lipgloss.NewStyle().Bold(true).Padding(0, 4).Border(lipgloss.Border{Bottom: lipgloss.NormalBorder().Bottom}),
			// Background color soft purple
			Selected: lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("#a448a4")),
			Cell:     lipgloss.NewStyle().Padding(0, 4).Border(lipgloss.Border{Top: "-", MiddleBottom: "-"}, true),
		}),
	)
}
