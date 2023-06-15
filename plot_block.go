package main

import "strings"

var (
	bchars = [...]rune{' ', '▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}
	Block  BlockPlot
)

type BlockPlot struct{}

func (BlockPlot) Display(vals []int, startIdx int, max int) string {
	var sb strings.Builder
	iter(vals, startIdx, func(v int) {
		sb.WriteRune(bchars[trns(v, max, 8)])
	})
	return sb.String()
}

func (BlockPlot) NewWidth(w int) int {
	return w
}
