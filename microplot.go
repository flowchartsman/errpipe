package main

import (
	"sync"
	"time"
)

type MicroplotConf struct {
	Width    int
	Max      int
	Interval time.Duration
	Style    Style
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
	adjWidth := conf.Style.NewWidth(conf.Width)
	m := &Microplot{
		max:      conf.Max,
		buckets:  make([]int, adjWidth),
		ticker:   time.NewTicker(conf.Interval),
		interval: conf.Interval,
		style:    conf.Style,
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
