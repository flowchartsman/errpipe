package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	signal.Ignore(syscall.SIGPIPE)
	config := appConfig{}
	flag.BoolVar(&config.displayWarning, "w", false, "display warnings")
	flag.BoolVar(&config.displayInfo, "i", false, "display info")
	flag.Parse()
	_, err := tea.NewProgram(
		newApp(config),
	).Run()
	if err != nil {
		fmt.Printf("Error running program:", err)
		os.Exit(1)
	}
}
