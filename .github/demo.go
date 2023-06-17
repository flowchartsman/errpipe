package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"time"
)

func printStuff(str string, times int, opt ...func()) {
	for i := 0; i < times; i++ {
		fmt.Println(str)
		for _, o := range opt {
			o()
		}
	}
}

func withDelay(millis time.Duration) func() {
	return func() {
		time.Sleep(millis * time.Millisecond)
	}
}

func ramp(width int, heights []int) {
	for _, height := range heights {
		printBlock(height, width)
	}
	for i := len(heights) - 2; i >= 0; i-- {
		printBlock(heights[i], width)
	}
}

func shortRamp(width int) {
	ramp(width, []int{1, 2, 3, 4})
}

func tallRamp(width int) {
	ramp(width, []int{1, 2, 3, 4, 5, 6, 7, 8})
}

func printBlock(height, width int) {
	for i := 0; i < width; i++ {
		printStuff(infom, height)
		idle(250)
	}
}

func idle(millis int) {
	time.Sleep(time.Duration(millis) * time.Millisecond)
}

const (
	warnm  = "[WARNING] This is a warning"
	errorm = "[ERROR] This is an error"
	infom  = "[INFO] This is an info"
)

func main() {
	var demoName string
	flag.StringVar(&demoName, "demo", "main", "which demo to play (default: \"long\"")
	flag.Parse()
	log.SetFlags(0)
	demo, found := demos[demoName]
	if !found {
		log.Println("invalid demo valid demos are:")
		names := make([]string, 0, len(demos))
		for k := range demos {
			names = append(names, k)
		}
		sort.Strings(names)
		for _, n := range names {
			log.Printf("  -%s", n)
		}
		os.Exit(1)
	}
	demo()
}

func demoMain() {
	shortRamp(3)
	time.Sleep(500 * time.Millisecond)
	printBlock(1, 3)
	printStuff(warnm, 1)
	printBlock(1, 2)
	printStuff(errorm, 1)
	printBlock(1, 1)
	printBlock(2, 1)
	printStuff(warnm, 2, withDelay(150))
	printBlock(2, 2)
}

func demoDelay() {
	demoMain()
	printBlock(3, 1)
	printBlock(1, 2)
	idle(7500)
}

func demoShort() {
	shortRamp(4)
	time.Sleep(10 * time.Second)
}

func demoTall() {
	tallRamp(2)
	time.Sleep(10 * time.Second)
}

func demoDouble() {
	shortRamp(4)
	time.Sleep(10 * time.Second)
}

var demos = map[string]func(){
	"main":   demoMain,
	"delay":  demoDelay,
	"short":  demoShort,
	"tall":   demoTall,
	"double": demoDouble,
}
