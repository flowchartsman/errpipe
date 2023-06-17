package main

import "strings"

type Braille struct {
	bar    bool
	fourUp bool
}

func NewBraille(barchart bool, fourUp bool) *Braille {
	return &Braille{
		bar:    barchart,
		fourUp: fourUp,
	}
}

func (b *Braille) Display(vals []int, startIdx int, max int) string {
	var sb strings.Builder
	left := true
	rnge := 3
	if b.fourUp {
		rnge = 4
	}
	var lv, rv int
	iter(vals, startIdx, func(v int) {
		if left {
			lv = trns(v, max, rnge)
			left = false
			return
		}
		rv = trns(v, max, rnge)
		if !b.fourUp {
			if lv > 0 {
				lv++
			}
			if rv > 0 {
				rv++
			}
		}
		sb.WriteRune(getchar(lv, rv, b.bar, b.fourUp))
		left = true
	})
	return sb.String()
}

func (b *Braille) NewWidth(w int) int {
	return w * 2
}

// 8 point braille rune layout
//
//	+------+
//	|(0)(3)|
//	|(1)(4)|
//	|(2)(5)|
//	|(6)(7)|
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

func getchar(l, r int, bar bool, fourUp bool) rune {
	bottom := 1
	if fourUp {
		bottom = 0
	}
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
			for i := l - 1; i > bottom; i-- {
				lowbits[lIdx(i)] = 1
			}
		}
	}

	if r > 0 {
		lowbits[rIdx(r)] = 1
		if bar {
			for i := r - 1; i > bottom; i-- {
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
