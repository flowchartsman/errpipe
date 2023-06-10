package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type appConfig struct {
	displayWarning bool
	displayInfo    bool
}

type appModel struct {
	messageQueue <-chan logMessage
	levels       map[logLevel]bool
	width        int
	height       int
	am           *activityMonitor
}

func newApp(config appConfig) appModel {
	return appModel{
		messageQueue: startMessageReader(),
		levels: map[logLevel]bool{
			lvlError: true,
			lvlWarn:  config.displayWarning,
			lvlInfo:  config.displayInfo,
		},
		am: NewActivityMonitor(20),
	}
}

func (m appModel) Init() tea.Cmd {
	return tea.Batch(
		m.nextMessage(),
		doUITick(),
	)
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	case logMessage:
		if msg.lvl == final {
			return m, tea.Sequence(
				tea.Println(m.am.Summarize()),
				tea.Quit,
			)
		}
		return m, tea.Batch(
			m.am.Measure(msg),
			m.Filter(msg),
			m.nextMessage(),
		)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
		return m, nil
	case uiTick:
		return m, doUITick()
	}
	return m, nil
}

func (m appModel) Filter(msg logMessage) tea.Cmd {
	if m.levels[msg.lvl] {
		return tea.Println(lipgloss.NewStyle().Width(m.width).Render(msg.msg))
	}
	return nil
}

func (m appModel) View() string {
	return m.am.View()
}

func (m appModel) nextMessage() tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-m.messageQueue
		if ok {
			return msg
		}
		return logMessage{lvl: final}
	}
}

type uiTick time.Time

func doUITick() tea.Cmd {
	return tea.Tick(150*time.Millisecond, func(t time.Time) tea.Msg {
		return uiTick(t)
	})
}
