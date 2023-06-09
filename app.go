package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type appConfig struct{}

type appModel struct {
	messages <-chan logMessage
	levels   map[logLevel]bool
	width    int
	height   int
	am       *activityMonitor
}

func newApp(levels map[logLevel]bool) appModel {
	return appModel{
		messages: newMessageReader(),
		levels:   levels,
		am:       NewActivityMonitor(10),
	}
}

func (m appModel) Init() tea.Cmd {
	return tea.Batch(
		m.nextMessage(),
		// m.am.Init(),
		doUITick(),
	)
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	case logMessage:
		return m, tea.Batch(
			m.am.Tick(),
			m.Filter(msg),
			m.nextMessage(),
		)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			// need a channel that can be closed and shutdown reader and waiter
			// or turn off signal handling
			// os.Exit(130)
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
		msg, ok := <-m.messages
		if ok {
			return msg
		}
		println("pipe complete")
		return tea.Quit
	}
}
