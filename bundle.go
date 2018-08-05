package main

import (
	"runtime"
	"time"
	"log"
	"syscall/js"
	"fmt"
)

const MaxCount = 60

func countUp(output chan int, ticker chan int) {
	var count = 0
	for {
		switch <- ticker {
		case 1:
			if count < MaxCount {
				count++
			} else {
				count = 0
			}
		case -1:
			if count > 0 {
				count--
			} else {
				count = MaxCount
			}
		case 0:
			// Do nothing this means that there is no update
		default:
			count = 0
		}
		output <- count
	}
}

func autoTick(ticker *time.Ticker, function func(args []js.Value)) {
	for range ticker.C {
		log.Println("RoboTick")
		function(nil)
	}
}

func main() {
	log.Printf("Hello Logger:\t%s - %s\n", runtime.GOOS, runtime.GOARCH)
	count := make(chan int)
	hidden := false
	forward := true
	ticker := make(chan int)

	bumpClick := func (args []js.Value) {
		if forward {
			ticker <- 1
		} else {
			ticker <- -1
		}
	}

	resetClick := func (args []js.Value) {
			ticker <- 2
	}

	reverseClick := func (args []js.Value) {
		forward = !forward
		ticker <- 0
	}

	globalDoc := js.Global().Get("document")
	display := globalDoc.Call("getElementById", "display")
	help := globalDoc.Call("getElementById", "help")
	bumpButton := globalDoc.Call("getElementById", "bumpButton")
	resetButton := globalDoc.Call("getElementById", "resetButton")
	reverseButton := globalDoc.Call("getElementById", "reverseButton")

	bump := js.NewCallback(bumpClick)
	reset := js.NewCallback(resetClick)
	reverse := js.NewCallback(reverseClick)

	log.Println("Adding callbacks")
	bumpButton.Call("addEventListener", "click", bump)
	resetButton.Call("addEventListener", "click", reset)
	reverseButton.Call("addEventListener", "click", reverse)

	log.Println("Starting ticker")
	roboTicker := time.NewTicker(5 * time.Second)
	go autoTick(roboTicker, bumpClick)

	log.Println("Starting counter")
	go countUp(count, ticker)

	for {
		// Get new count
		currentCount := <- count

		// Update Text
		output := fmt.Sprintf("Counter: %d", currentCount)
		display.Set("innerHTML", output)

		// Toggle Style class
		reverseClasses := reverseButton.Get("classList")
		if !forward {
			reverseClasses.Call("add", "reverse")
		} else {
			reverseClasses.Call("remove", "reverse")
		}

		// Set Text Color based on count
		bumpScale := int(float64(currentCount)/ float64(MaxCount) * 255)
		bumpColor := fmt.Sprintf("color: #%02x%02x00;", bumpScale, 255-bumpScale)
		display.Call("setAttribute", "style", bumpColor)

		// Hide help text
		if !hidden {
			hidden = true
			help.Call("setAttribute", "hidden", "")
		}
	}
}
