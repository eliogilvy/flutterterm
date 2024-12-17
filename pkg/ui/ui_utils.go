package ui

import (
	"flutterterm/pkg/utils"

	"github.com/charmbracelet/bubbles/spinner"
)

type devicesComplete []utils.Device

type runningComplete bool

type cmdError string

type state int

const (
	// Viewing stuff
	view state = iota
	// Loading stuff
	getting
	// Running command
	running
)

// Default loading spinner
func getSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return s
}

const quitAndHelpMessage = "\nPress q to quit, or ? for help\n"
const controlsHelpMessage = "Controls\nj, down: go down\nk, up: go up\nh, left: go back (if applicable)\nenter: submit\n\n"
