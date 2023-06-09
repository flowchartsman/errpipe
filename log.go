package main

type logLevel int

const (
	lvlNone logLevel = iota
	lvlInfo
	lvlWarn
	lvlError
)

type logMessage struct {
	lvl logLevel
	msg string
}
