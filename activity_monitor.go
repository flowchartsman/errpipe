package main

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type uiTick time.Time

func doUITick() tea.Cmd {
	return tea.Tick(150*time.Millisecond, func(t time.Time) tea.Msg {
		return uiTick(t)
	})
}

type activityMonitor struct {
	idle    bool
	width   int
	lastMsg time.Time
	plot    *Microplot
}

func NewActivityMonitor(width int) *activityMonitor {
	return &activityMonitor{
		width:   width,
		lastMsg: time.Now(),
		plot: NewMicroplot(MicroplotConf{
			Width:    width - 2,
			Max:      4,
			Interval: 250 * time.Millisecond,
			Style:    NewBraille(true),
		}),
	}
}

func (am *activityMonitor) Init() tea.Cmd {
	return nil
}

func (am *activityMonitor) View() string {
	var sb strings.Builder
	since := time.Since(am.lastMsg)
	if !am.idle && since > 2*time.Second {
		am.idle = true
		am.plot.Pause()
	}
	sb.WriteString("[")
	sb.WriteString(am.plot.String())
	sb.WriteString("]")
	if am.idle {
		sb.WriteString("⏳")
		if since > 5*time.Second {
			sb.WriteString(" " + fmtDuration(since))
		}
	}

	return sb.String()
}

func (am *activityMonitor) Update(msg tea.Msg) (*activityMonitor, tea.Cmd) {
	return am, nil
}

// update better?
func (am *activityMonitor) Tick() tea.Cmd {
	am.lastMsg = time.Now()
	am.idle = false
	am.plot.Measure(1)
	return nil
}
