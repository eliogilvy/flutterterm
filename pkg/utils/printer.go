package utils

import (
	"github.com/fatih/color"
)

func PrintError(s string) {
	color.Red(s)
}

func PrintHelp(s string) {
	color.Yellow(s)
}

func PrintInfo(s string) {
	color.Blue(s)
}

func PrintSuccess(s string) {
	color.Green(s)
}

func PrintWarning(s string) {
	color.Yellow(s)
}
