package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	config := appConfig{}
	flag.BoolVar(&config.displayWarning, "w", false, "display warnings")
	flag.BoolVar(&config.displayInfo, "i", false, "display info")
	flag.Parse()
	if _, err := tea.NewProgram(newApp(config)).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
