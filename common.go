package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

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
	if hours == 0 {
		sb.WriteString(fmt.Sprintf(".%d", milli))
	}
	sb.WriteString("s")

	return sb.String()
}

func trns(v, oldmax, newmax int) int {
	return v * newmax / oldmax
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
