package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/charmbracelet/lipgloss"
)

var (
	rxHasPfx = regexp.MustCompile(`(?i)^\s*[ \[(|]?(err(or)?|warn(ing)?|info?)[ \])|]`)
	rxError  = regexp.MustCompile(`(?i)err(or)?`)
	rxWarn   = regexp.MustCompile(`(?i)warn(ing)?`)
	rxInfo   = regexp.MustCompile(`(?i)info?`)
	rxLeadSp = regexp.MustCompile(`^\s+`)
)

func isContinuation(rawMsg string) bool {
	return rawMsg == "" || rxLeadSp.MatchString(rawMsg)
}

func colorize(msg string, loc []int, lvl logLevel) string {
	var (
		style lipgloss.Style
		sb    strings.Builder
	)
	switch lvl {
	case lvlError:
		style = errStyle
	case lvlWarn:
		style = warnStyle
	case lvlInfo:
		style = infoStyle
	}
	sb.WriteString(msg[:loc[0]])
	sb.WriteString(style.Render(msg[loc[0]:loc[1]]))
	sb.WriteString(msg[loc[1]:])
	return sb.String()
}

func getLeveledMsg(msg string) logMessage {
	lvl := lvlNone
	pfxLoc := rxHasPfx.FindStringIndex(msg)
	if len(pfxLoc) > 0 {
		pfx := msg[pfxLoc[0]:pfxLoc[1]]
		switch {
		case rxError.MatchString(pfx):
			lvl = lvlError
		case rxWarn.MatchString(pfx):
			lvl = lvlWarn
		case rxInfo.MatchString(pfx):
			lvl = lvlInfo
		}
		msg = colorize(msg, pfxLoc, lvl)
	}
	return logMessage{lvl, msg}
}

var (
	errStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#D24334"))
	warnStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#EABA4C"))
	infoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#123B68"))
)

// probably want err return
func newMessageReader() <-chan logMessage {
	out := make(chan logMessage, 100)
	go func() {
		lastLvl := lvlNone
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			rawMsg := stripansi.Strip(sc.Text())
			// if isContinuation(rawMsg) {
			// 	out <- logMessage{lastLvl, rawMsg}
			// 	continue
			// }
			msg := getLeveledMsg(rawMsg)
			if msg.lvl == lvlNone {
				msg.lvl = lastLvl
			}
			lastLvl = msg.lvl
			out <- msg
		}
		close(out)
		if sc.Err() != nil {
			fmt.Fprintf(os.Stderr, "err while reading: %v", sc.Err())
		}
	}()
	return out
}
