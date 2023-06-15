package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/acarl005/stripansi"
)

var (
	rxHasPfx = regexp.MustCompile(`(?i)^\s*[\[(|]?(err(or)?|warn(ing)?|info?)[ \])|]`)
	rxError  = regexp.MustCompile(`(?i)err(or)?`)
	rxWarn   = regexp.MustCompile(`(?i)warn(ing)?`)
	rxInfo   = regexp.MustCompile(`(?i)info?`)
	rxSpc    = regexp.MustCompile(`^s*$`)
)

func colorize(msg string, loc []int, lvl logLevel) string {
	var sb strings.Builder
	sb.WriteString(msg[:loc[0]])
	sb.WriteString(lvl.format(msg[loc[0]:loc[1]]))
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
	return logMessage{
		lvl: lvl,
		msg: msg,
	}
}

func startMessageReader() <-chan logMessage {
	out := make(chan logMessage)
	go func() {
		lastLvl := lvlNone
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			rawMsg := stripansi.Strip(sc.Text())
			logmsg := getLeveledMsg(rawMsg)
			// empty lines reset the level
			if rxSpc.MatchString(logmsg.msg) {
				logmsg.continuation = true
				logmsg.lvl = lvlNone
			} else {
				// if message doesn't have level, but last message was leveled, assume continuation
				if logmsg.lvl == lvlNone && lastLvl != lvlNone {
					logmsg.lvl = lastLvl
				}
			}
			lastLvl = logmsg.lvl
			out <- logmsg
		}
		close(out)
		if sc.Err() != nil {
			fmt.Fprintf(os.Stderr, "err while reading: %v", sc.Err())
		}
	}()
	return out
}
