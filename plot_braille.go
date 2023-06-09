package main

import "strings"

type Braille struct {
	bar bool
}

func NewBraille(barchart bool) *Braille {
	return &Braille{
		bar: barchart,
	}
}

func (b *Braille) Display(vals []int, startIdx int, max int) string {
	var sb strings.Builder
	left := true
	var lv, rv int
	iter(vals, startIdx, func(v int) {
		if left {
			lv = trns(v, max, 4)
			left = false
			return
		}
		rv = trns(v, max, 4)
		sb.WriteRune(getchar(lv, rv, b.bar))
		left = true
	})
	return sb.String()
}

// 8 point braille rune layout
//
//	+------+
//	|(1)(4)|
//	|(2)(5)|
//	|(3)(6)|
//	|(7)(8)|
//	+------+
//
// See https://en.wikipedia.org/wiki/Braille_Patterns#Identifying.2C_naming_and_ordering)

func lIdx(n int) int {
	if n == 1 {
		return 6
	}
	return 6 - 2 - n
}

func rIdx(n int) int {
	if n == 1 {
		return 7
	}
	return lIdx(n) + 3
}

func getchar(l, r int, bar bool) rune {
	lowbits := [8]int{}
	if l > 4 {
		l = 4
	}
	if r > 4 {
		r = 4
	}
	if l > 0 {
		lowbits[lIdx(l)] = 1
		if bar {
			for i := l - 1; i > 0; i-- {
				lowbits[lIdx(i)] = 1
			}
		}
	}

	if r > 0 {
		lowbits[rIdx(r)] = 1
		if bar {
			for i := r - 1; i > 0; i-- {
				lowbits[rIdx(i)] = 1
			}
		}
	}

	var v int
	for i, x := range lowbits {
		v += x << uint(i)
	}
	return rune(v) + '\u2800'
}
