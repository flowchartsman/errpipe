package main

import (
	"sync"
	"time"
)

type MicroplotConf struct {
	Width    int
	Max      int
	Interval time.Duration
	Style    string
}

type Microplot struct {
	mu       sync.RWMutex
	closer   sync.Once
	closed   bool
	paused   bool
	ticker   *time.Ticker
	interval time.Duration
	max      int
	buckets  []int
	write    int
	style    Style
}

func NewMicroplot(conf MicroplotConf) *Microplot {
	var style Style
	switch conf.Style {
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
	adjWidth := style.NewWidth(conf.Width)
	m := &Microplot{
		max:      conf.Max,
		buckets:  make([]int, adjWidth),
		ticker:   time.NewTicker(conf.Interval),
		interval: conf.Interval,
		style:    style,
	}
	go func() {
		for range m.ticker.C {
			m.mu.Lock()
			if !m.paused {
				m.shift()
			}
			m.mu.Unlock()
		}
	}()
	return m
}

func (m *Microplot) Pause() {
	m.mu.Lock()
	m.paused = true
	m.mu.Unlock()
}

func (m *Microplot) Close() {
	m.closer.Do(func() {
		m.mu.Lock()
		m.ticker.Stop()
		m.closed = true
		m.mu.Unlock()
	})
}

func (m *Microplot) Measure(v int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return
	}
	if m.paused {
		m.shift()
		m.paused = false
		m.ticker.Reset(m.interval)
	}
	m.buckets[m.write] += v
	if m.buckets[m.write] > m.max {
		m.buckets[m.write] = m.max
	}
}

func (m *Microplot) shift() {
	m.write--
	if m.write == -1 {
		m.write = len(m.buckets) - 1
	}
	m.buckets[m.write] = 0
}

func (m *Microplot) String() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.style.Display(m.buckets, m.write, m.max)
}

type Style interface {
	Display(vals []int, startIdx int, max int) string
	NewWidth(int) int
}
