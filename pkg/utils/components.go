package utils

import (
	"github.com/charmbracelet/bubbles/list"
)

// Return styled list
func GetList(items []list.Item) list.Model {
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)

	l.KeyMap.CursorDown.SetEnabled(true)
	l.KeyMap.CursorUp.SetEnabled(true)
	l.InfiniteScrolling = true

	l.SetSize(100, len(items)*7)

	return l
}
