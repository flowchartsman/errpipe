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
	if conf.Width <= 0 {
		conf.Width = 10
	}
	if conf.Interval <= 0 {
		conf.Interval = 250 * time.Millisecond
	}
	if conf.Max < 4 {
		conf.Max = 4
	}
	if _, ok := conf.Style.(*Braille); ok {
		conf.Width *= 2
	}
	m := &Microplot{
		max:      conf.Max,
		buckets:  make([]int, conf.Width),
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
}
