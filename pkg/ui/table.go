package ui

import (
	"flutterterm/pkg/model"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

type TableModel = *table.Model
type TableRow = table.Row
type TableColumn = table.Column

func GetTable(c []TableColumn, r []TableRow) TableModel {
	t := table.New(
		table.WithColumns(c),
		table.WithRows(r),
		table.WithHeight(len(r)+1),
		table.WithStyles(table.Styles{
			Header:   lipgloss.NewStyle().Padding(1, 0),
			Cell:     lipgloss.Style{},
			Selected: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")),
		}),
	)
	return &t
}

func GetDeviceTable(devices []model.Device) TableModel {
	c := []TableColumn{
		{Title: "Name", Width: 40},
		{Title: "ID", Width: 40},
	}

	var r []TableRow

	for _, device := range devices {
		row := TableRow{
			device.Name,
			device.ID,
		}

		r = append(r, row)
	}

	t := GetTable(c, r)

	return t
}

func GetConfigTable(configs []model.FlutterConfig) TableModel {
	c := []TableColumn{
		{Title: "Config", Width: 40},
		{Title: "Description", Width: 40},
	}

	var r []TableRow

	for _, config := range configs {
		row := TableRow{config.Name, config.Description}

		r = append(r, row)
	}

	t := GetTable(c, r)

	return t
}
