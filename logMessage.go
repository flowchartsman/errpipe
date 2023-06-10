package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type logLevel int

const (
	lvlNone logLevel = iota
	lvlInfo
	lvlWarn
	lvlError
	final
)

func (l logLevel) short() string {
	switch l {
	case lvlInfo:
		return "I"
	case lvlWarn:
		return "W"
	case lvlError:
		return "E"
	}
	return "?"
}

func (l logLevel) summary() string {
	switch l {
	case lvlInfo:
		return "Info"
	case lvlWarn:
		return "Warnings"
	case lvlError:
		return "Errors"
	}
	return "?"
}

func (l logLevel) format(s string, a ...any) string {
	style := stylNone
	switch l {
	case lvlInfo:
		style = stylInfo
	case lvlWarn:
		style = stylWarn
	case lvlError:
		style = stylErr
	}
	return style.Render(fmt.Sprintf(s, a...))
}

var (
	stylErr  = lipgloss.NewStyle().Foreground(lipgloss.Color("#D24334"))
	stylWarn = lipgloss.NewStyle().Foreground(lipgloss.Color("#EABA4C"))
	stylInfo = lipgloss.NewStyle().Foreground(lipgloss.Color("#65816A"))
	stylNone = lipgloss.NewStyle().Foreground(lipgloss.NoColor{})
)

type logMessage struct {
	lvl          logLevel
	msg          string
	continuation bool
}
