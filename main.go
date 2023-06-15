package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	config := appConfig{}
	flag.BoolVar(&config.displayWarning, "w", false, "display warnings")
	flag.BoolVar(&config.displayInfo, "i", false, "display info")
	flag.Parse()
	config.max = ienvd("ERRPIPE_MAX", 3)
	config.width = ienvd("ERRPIPE_WIDTH", 20)
	config.style = envd("ERRPIPE_STYLE", "braille")
	config.intervalMs = ienvd("ERRPIPE_INTERVAL", 250)
	_, err := tea.NewProgram(
		newApp(config),
	).Run()
	if err != nil {
		fmt.Printf("Errpipe: %v", err)
		os.Exit(1)
	}
}

func envd(varName string, defaultValue string) string {
	envStr, found := os.LookupEnv(varName)
	if !found {
		return defaultValue
	}
	return envStr
}

func ienvd(varName string, defaultValue int) int {
	envStr, found := os.LookupEnv(varName)
	if !found {
		return defaultValue
	}
	iv, err := strconv.Atoi(envStr)
	if err != nil {
		return defaultValue
	}
	return iv
}
