package main

import (
	"fmt"

	"github.com/kashisau/cec-cascade/devices"
)

type ActiveSource int

func main() {
	// Initialisation
	tvOn := false
	soundOn := false

	// Sound state channels
	soundStatus := make(chan bool)
	go devices.WatchSoundOutput(soundStatus)

	tvStatus := make(chan bool)
	go devices.WatchHdmi(tvStatus)

	// Check the channel values until termination
	for {
		select {
		case soundStatusUpdate := <-soundStatus:
			soundOn = soundStatusUpdate
		case tvStatusUpdate := <-tvStatus:
			tvOn = tvStatusUpdate
		}

		fmt.Printf("TV: %t Sound: %t\n", tvOn, soundOn)
	}
}
