package ui

import (
	"flutterterm/pkg/utils"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type EmulatorModel struct {
	devices          []utils.Device
	selectedEmulator utils.Device
	state            state
	spinner          spinner.Model
	isCold           bool // Cold start
	list             list.Model
}

func InitialEmulatorModel(isCold bool) EmulatorModel {
	return EmulatorModel{
		state:   getting,
		spinner: getSpinner(),
		isCold:  isCold,
	}
}

func (m EmulatorModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, getEmulators())
}

func (m EmulatorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "?":
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			m.list.CursorUp()
			return m, nil
		case "down", "j":
			m.list.CursorDown()
			return m, nil
		case "enter":
			m.selectedEmulator = m.devices[m.list.Index()]
			m.state = running
			return m, launchEmulator(m)
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case devicesComplete:
		m.state = view
		m.devices = msg
		m.list = utils.GetDeviceList(m.devices)
		return m, nil
	case runningComplete:
		return m, tea.Quit
	case cmdError:
		utils.PrintError(fmt.Sprintf("%s", msg))
		m.state = view
		return m, tea.Quit
	}
	return m, nil
}

func (m EmulatorModel) View() string {
	switch m.state {
	case view:
		return DocStyle.Render(m.list.View())
	case getting:
		spinner := m.spinner.View()
		return fmt.Sprintf("%s Getting emulators %s", spinner, spinner)
	case running:
		spinner := m.spinner.View()
		return fmt.Sprintf("%s Launching %s %s", spinner, m.selectedEmulator.Name, spinner)
	default:
		return "Unknown state"
	}

}

func getEmulators() tea.Cmd {
	return func() tea.Msg {
		cmd := utils.FlutterEmulators([]string{})

		output, err := cmd.Output()

		if err != nil {
			return cmdError(err.Error())
		}

		devices, err := utils.ParseEmulators(output)
		if err != nil {
			return cmdError(err.Error())
		}

		return devicesComplete(devices)
	}
}

func launchEmulator(m EmulatorModel) tea.Cmd {
	return func() tea.Msg {
		isCold := m.isCold
		args := []string{"--launch", m.selectedEmulator.ID}

		if isCold {
			args = append(args, "--cold")
		}

		cmd := utils.FlutterEmulators(args)
		err := cmd.Run()
		if err != nil {
			return cmdError(err.Error())
		}
		return runningComplete(true)
	}
}
