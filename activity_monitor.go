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

func NewActivityMonitor(config appConfig) *activityMonitor {
	var style Style
	switch config.style {
	case "braille-line":
		style = NewBraille(false, false)
	case "braille4-line":
		style = NewBraille(false, true)
	case "braille4":
		style = NewBraille(true, true)
	case "block":
		style = Block
	case "legacy":
		style = TwoTuplePlot{LegacyLine}
	case "legacy-block":
		style = TwoTuplePlot{LegacyBlock}
	case "legacy-block-line":
		style = TwoTuplePlot{LegacyBlockLine}
	default:
		style = NewBraille(true, false)
	}
	return &activityMonitor{
		width:   config.width,
		lastMsg: time.Now(),
		plot: NewMicroplot(MicroplotConf{
			Width:    config.width - 2, // for the brackets
			Max:      20,
			Interval: 250 * time.Millisecond,
			Style:    style,
		}),
		counts: map[logLevel]int{},
	}
}

func (*activityMonitor) Init() tea.Cmd {
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
	if am.idle && since > 5*time.Second {
		sb.WriteString("⏳")
		sb.WriteString(" " + fmtDuration(since))
		sb.WriteString(strings.Repeat(" ", am.width-sb.Len()))
	} else {
		sb.WriteString(am.plot.String())
	}
	sb.WriteString("] ")
	for lvl := lvlError; lvl >= lvlNone; lvl-- {
		// dont print "none" messages until they've occurred
		if lvl == lvlNone && am.counts[lvl] == 0 {
			continue
		}
		sb.WriteString(lvl.format("%s", lvl.short()))
		sb.WriteString(fmt.Sprintf(":%d ", am.counts[lvl]))
	}
	return sb.String()
}

func (am *activityMonitor) Summarize() string {
	var sb strings.Builder
	sb.WriteString("\n")
	for lvl := lvlError; lvl >= lvlNone; lvl-- {
		// dont print "none" messages unless they've occurred
		if lvl == lvlNone && am.counts[lvl] == 0 {
			continue
		}
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
