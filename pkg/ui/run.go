package ui

import (
	"errors"
	"flutterterm/pkg/utils"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type RunModel struct {
	devices         []utils.Device
	configs         []utils.FlutterRunConfig
	stage           devicestage
	Selected_device utils.Device
	Selected_config utils.FlutterRunConfig
	state           state
	spinner         spinner.Model
	deviceList      list.Model
	configList      list.Model
}

type devicestage int

const (
	device devicestage = iota
	config
	_length
)

func InitialRunModel(configs []utils.FlutterRunConfig) RunModel {
	return RunModel{
		configs:    configs,
		stage:      device,
		state:      getting,
		spinner:    getSpinner(),
		configList: utils.GetConfigList(configs),
	}
}

func (m RunModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, getDevices())
}

func (m RunModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		h, v := DocStyle.GetFrameSize()
		if m.stage == device && m.state == view {
			m.deviceList.SetSize(msg.Width-h, msg.Height-v)
		}
		m.configList.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:

		switch msg.String() {

		case "?":
			// TODO

		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			m = m.goUp()
		case "down", "j":
			m = m.goDown()
		case "left", "h":
			m, _ := m.back()
			return m, nil
		case "enter":
			m, cmd := m.doNextThing()
			return m, cmd
		}
		return m, nil

	case DeviceSelectedMsg:
	case devicesComplete:
		m.devices = msg
		m.state = view
		m.deviceList = utils.GetDeviceList(m.devices)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m RunModel) goUp() RunModel {
	switch m.stage {
	case config:
		m.configList.CursorUp()
	case device:
		m.deviceList.CursorUp()
	default: // Do nothng
	}
	return m
}

func (m RunModel) goDown() RunModel {
	switch m.stage {
	case config:
		m.configList.CursorDown()
	case device:
		m.deviceList.CursorDown()
	default: // Do nothng
	}
	return m
}

// Go back in the process
func (m RunModel) back() (RunModel, error) {
	if m.stage == device {
		return m, errors.New("Already at beginning")
	}
	m.stage = device
	return m, nil
}

// Go to the next part of the process
func (m RunModel) doNextThing() (RunModel, tea.Cmd) {
	var cmd tea.Cmd
	switch m.stage {
	case device:
		m.Selected_device = m.devices[m.deviceList.Index()]
		m.stage = config
		cmd = DeviceSelected
	case config:
		m.Selected_config = m.configs[m.configList.Index()]
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
	switch m.state {
	case view:
		switch m.stage {
		case device:
			s += DocStyle.Render(m.deviceList.View())
		case config:
			s += fmt.Sprintf("Device: %s\n\n", m.Selected_device.Name)
			s += DocStyle.Render(m.configList.View())
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
