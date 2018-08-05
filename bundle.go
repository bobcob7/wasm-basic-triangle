package main

import (
	"runtime"
	"log"
	"syscall/js"
	"time"
	"fmt"
)

type IPInfo struct {
	City string `json:"city"`
	Country string `json:"country"`
	Hostname string `json:"hostname"`
	IP string `json:"ip"`
	IPDecimal string `json:"ip_decimal"`
}

var counter = 0

func countUp(counter chan int) {
	current := 0
	for {
		if current < 5000 {
			current++
		} else {
			current = 0
		}
		counter <- current
		time.Sleep(time.Second * 1)
	}
}

func main() {
	log.Printf("Hello Logger:\t%s - %s\n", runtime.GOOS, runtime.GOARCH)
	counter := make(chan int)

	globalDoc := js.Global().Get("document")

	el := globalDoc.Call("getElementById", "thing")

	go countUp(counter)

	for {
		count := <- counter
		log.Printf("Count: %d\n", count)
		if count % 5 == 0 {
			// Update count
			text := fmt.Sprintf("Current Count: %d", count)
			el.Set("innerHTML", text)
		}
	}
}
