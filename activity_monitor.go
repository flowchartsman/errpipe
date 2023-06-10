package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type activityMonitor struct {
	idle    bool
	width   int
	lastMsg time.Time
	plot    *Microplot
	counts  map[logLevel]int
}

func NewActivityMonitor(width int) *activityMonitor {
	return &activityMonitor{
		width:   width,
		lastMsg: time.Now(),
		plot: NewMicroplot(MicroplotConf{
			Width:    width - 2,
			Max:      20,
			Interval: 250 * time.Millisecond,
			Style:    NewBraille(true),
		}),
		counts: map[logLevel]int{},
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
	for lvl := lvlError; lvl > lvlNone; lvl-- {
		sb.WriteString(lvl.format("%s", lvl.short()))
		sb.WriteString(fmt.Sprintf(":%d ", am.counts[lvl]))
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

func (am *activityMonitor) Summarize() string {
	var sb strings.Builder
	sb.WriteString("\n")
	for lvl := lvlError; lvl > lvlNone; lvl-- {
		sb.WriteString(lvl.format("%s", lvl.short()))
		sb.WriteString(fmt.Sprintf(":%d ", am.counts[lvl]))
	}
	return sb.String()
}

func (am *activityMonitor) Update(msg tea.Msg) (*activityMonitor, tea.Cmd) {
	return am, nil
}

func (am *activityMonitor) Measure(msg logMessage) tea.Cmd {
	am.lastMsg = time.Now()
	if msg.continuation {
		return nil
	}
	am.counts[msg.lvl]++
	am.idle = false
	am.plot.Measure(1)
	return nil
}

func fmtDuration(duration time.Duration) string {
	hours := int64(math.Mod(duration.Hours(), 24))
	minutes := int64(math.Mod(duration.Minutes(), 60))
	seconds := int64(math.Mod(duration.Seconds(), 60))
	milli := duration.Milliseconds() % 1000 / 100

	var sb strings.Builder
	if hours > 0 {
		sb.WriteString(fmt.Sprintf("%dh ", hours))
	}
	if minutes > 0 || hours > 0 {
		sb.WriteString(fmt.Sprintf("%dm ", minutes))
	}
	sb.WriteString(fmt.Sprintf("%d", seconds))
	if minutes == 0 {
		sb.WriteString(fmt.Sprintf(".%d", milli))
	}
	sb.WriteString("s")

	return sb.String()
}

func iter(buckets []int, start int, f func(v int)) {
	c := start
	for {
		f(buckets[c])
		c++
		if c == len(buckets) {
			c = 0
		}
		if c == start {
			break
		}
	}
}

func trns(v, oldmax, newmax int) int {
	nv := v * newmax / oldmax
	if v > 0 && nv == 0 {
		return 1
	}
	return nv
}
