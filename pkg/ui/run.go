package ui

import (
	"errors"
	"flutterterm/pkg/utils"
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	// "github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type RunModel struct {
	devices         []utils.Device
	configs         []utils.FlutterRunConfig
	cursor          utils.Navigator
	stage           devicestage
	Selected_device utils.Device
	Selected_config utils.FlutterRunConfig
	state           state
	spinner         spinner.Model
	list            list.Model
}

type devicestage int

const (
	device devicestage = iota
	config
	_length
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func InitialRunModel(configs []utils.FlutterRunConfig) RunModel {
	return RunModel{
		configs: configs,
		stage:   device,
		state:   getting,
		spinner: getSpinner(),
	}
}

func (m RunModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, getDevices())
}

func (m RunModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		if m.stage == device && m.state == view {
			h, v := docStyle.GetFrameSize()
			m.list.SetSize(msg.Width-h, msg.Height-v)
		}

	case tea.KeyMsg:

		switch msg.String() {

		case "?":
			m.cursor.ToggleHelp()

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			m.cursor.Previous()
			m.list.CursorUp()
		case "down", "j":
			m.cursor.Next()
			m.list.CursorDown()
		case "left", "h":
			m, err := m.back()
			if err == nil {
				m.cursor = utils.NewNavigator(0, len(m.devices))
			}
			return m, nil
		case "enter":
			m, cmd := m.doNextThing()
			return m, cmd
		}
		return m, nil

	case devicesComplete:
		m.devices = msg
		m.cursor = utils.NewNavigator(0, len(m.devices))
		m.state = view
		m.list = utils.GetDeviceList(m.devices)
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

// Go back in the process
func (m RunModel) back() (RunModel, error) {
	if m.stage == device {
		return m, errors.New("Couldn't go back")
	}
	m.stage = device
	return m, nil
}

// Go to the next part of the process
func (m RunModel) doNextThing() (RunModel, tea.Cmd) {
	var cmd tea.Cmd
	switch m.stage {
	case device:
		m.Selected_device = m.devices[m.list.Index()]
		m.cursor.Reset(len(m.configs))
		m.stage = config
		cmd = nil
	case config:
		m.Selected_config = m.configs[m.cursor.Index()]
		cmd = tea.Quit
	}
	return m, cmd
}

// Whether the model has enough information to run
func (m RunModel) IsComplete() bool {
	return m.Selected_config.Name != "" && m.Selected_device.ID != ""
}

func (m RunModel) View() string {
	var s string = ""
	if m.cursor.ShouldShowHelp() {
		s += controlsHelpMessage
	}
	switch m.state {
	case view:
		switch m.stage {
		case device:
			return docStyle.Render(m.list.View())
		case config:
			s += fmt.Sprintf("Device: %s\n\n", m.Selected_device.Name)
			s += "Select a config\n\n"
			for i, config := range m.configs {
				cursor := " "
				if m.cursor.Index() == i {
					cursor = utils.CursorChar
				}
				s += fmt.Sprintf("%s %s\n", cursor, config.Name)
			}
		}
		return s
	case getting:
		spinner := m.spinner.View()
		s := fmt.Sprintf("%s Getting devices %s", spinner, spinner)
		return s
	default:
		return "Unknown state"
	}
}

func getDevices() tea.Cmd {
	return func() tea.Msg {
		cmd := utils.FlutterDevices()
		output, err := cmd.Output()

		if err != nil {
			return cmdError(err.Error())
		}

		devices, err := utils.ParseDevices(output)

		if err != nil {
			return cmdError(err.Error())
		}

		return devicesComplete(devices)
	}
}
