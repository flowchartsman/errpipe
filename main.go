package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	logLevels := map[logLevel]bool{lvlError: true}
	flag.Func("w", "include warnings", lvlFlag(lvlWarn, logLevels))
	flag.Func("i", "include info", lvlFlag(lvlWarn, logLevels))
	flag.Parse()
	if _, err := tea.NewProgram(newApp(logLevels)).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func lvlFlag(lvl logLevel, m map[logLevel]bool) func(string) error {
	return func(sv string) error {
		if sv != "" {
			return fmt.Errorf("flag does not take an argument")
		}
		m[lvl] = true
		return nil
	}
}
