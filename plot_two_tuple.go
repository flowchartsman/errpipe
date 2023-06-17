package main

import "strings"

type TwoTuplePlotStyle int

const (
	LegacyLine TwoTuplePlotStyle = iota
	LegacyBlock
	LegacyBlockLine
)

var twoTupleChars = [...][4][4]rune{
	LegacyLine: {
		{' ', '🭈', '🭊', '🭋'},
		{'🬽', '🬭', '🭆', '🭄'},
		{'🬿', '🭑', '🬹', '🭂'},
		{'🭀', '🭏', '🭍', '🮋'},
	}, //               ⬆ sub
	LegacyBlock: { /*    ⬇ sub */
		{' ', '🬞', '🬦', '▐'},
		{'🬏', '🬭', '🬵', '🬷'},
		{'🬓', '🬱', '🬹', '🬻'},
		{'▌', '🬲', '🬺', '🮋'},
	}, /*⬆ sub */
	LegacyBlockLine: {
		{' ', '🬞', '🬦', '🬘'},
		{'🬏', '🬭', '🬖', '🬔'},
		{'🬃', '🬢', '🬋', '🬅'},
		{'🬣', '🬧', '🬈', '🬂'},
	},
}

type TwoTuplePlot struct {
	style TwoTuplePlotStyle
}

func (t TwoTuplePlot) Display(vals []int, startIdx int, max int) string {
	var sb strings.Builder
	last := trns(vals[startIdx], max, 3)
	first := true
	iter(vals, startIdx, func(v int) {
		if first {
			first = false
			return
		}
		v = trns(v, max, 3)
		if !first {
			sb.WriteRune(twoTupleChars[t.style][last][v])
		}
		last = v
	})
	return sb.String()
}

func (TwoTuplePlot) NewWidth(w int) int {
	return w + 1
}

/*
LEGACY SEXTANTS:
🬀 🬁 🬂 🬃 🬄 🬅 🬆 🬇 🬈 🬉 🬊 🬋 🬌 🬍 🬎 🬏 🬐 🬑 🬒 🬓  🬔 🬕 🬖 🬗 🬘 🬙 🬚 🬛 🬜 🬝 🬞 🬟 🬠 🬡 🬢 🬣 🬤 🬥 🬦 🬧 🬨 🬩 🬪 🬫 🬬 🬭 🬮 🬯 🬰 🬱 🬲 🬳 🬴 🬵 🬶 🬷 🬸 🬹 🬺 🬻
MISSING SEXTANTS:
BLOCK SEXTANT-135    - replacement:'▌'
BLOCK SEXTANT-246    - replacement:'▐'
BLOCK SEXTANT-123456 - replacement: '█' FULL BLOCK or '🮋' LEGACY LEFT 3/4 Block (close enough)
None of these are really good though.
Why can't the Unicode Consortium just be thorough?

LEGACY LINEGRAPH CHARS:
🬼 🬽 🬾 🬿 🭀 🭁 🭂 🭃 🭄 🭅 🭆 🭇 🭈 🭉 🭊 🭋 🭌 🭍 🭎 🭏
🭐 🭑 🭒 🭓 🭔 🭕 🭖 🭗 🭘 🭙 🭚 🭛 🭜 🭝 🭞 🭟
🭠 🭡 🭢 🭣 🭤 🭥 🭦 🭧
*/

/*
TODO: smoothing code instead of 0<->3 0<->2 and 1<->3 transitions so that
0->3->0->0 goes from
`🬘🬣 `
to:
`🬖🬈🬏`
0->3->0->1
`🬘🬣🬞`
to:
`🬖🬈🬭`
*/
