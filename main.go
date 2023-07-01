package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kashisau/cec-cascade/devices"
)

type ActiveSource int

const SOURCE_SOUND_OUTPUT = "260048001c1c3a3b3a1e1c1e1d1e1c1e1d3b1c1e3a1e1c1e1d000b971d1d3a3b3b1d1c1e1d1e1c1e1d3a1d1e3a1d1d1e1d000b971d1d3a3b3a1e1c1e1d1d1d1e1d3a1d1e3a1d1d1e1d000d05"
const SOURCE_TV = "260034001d1b1f1c1f1b3c1c1e1c1f1b1f1c1f381f1c1f1b3c1c1e000b941f1b1f1c1e1c3c1b1f1c1f1b1f1c1f381f1b1f1c3c1b1f000d05"
const SOURCE_AIRPLAY = "260030001b1d3a3b3a1e1c1e1d1d1d1e1d3a1d1e1c1e3a3b1c000b791d1d3a3b3a1d1d1e1d1d1d1e1c3b1d1d1d1e3a3a1d000d05"

const ON_SYMBOL = "✅"
const OFF_SYMBOL = "❌"

// Initialisation
var tvOn = false
var soundOuputOn = false

var tvStateChar = OFF_SYMBOL
var soundOutputChar = OFF_SYMBOL
var airplayStateChar = OFF_SYMBOL

func main() {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	WEATHER_URL := os.Getenv("WEATHER_URL")
	ROON_DISPLAY_URL := os.Getenv("ROON_DISPLAY_URL")

	// Sound state channels
	soundOutputStatus := make(chan bool)
	go devices.WatchSoundOutput(soundOutputStatus)

	tvStatus := make(chan bool)
	go devices.WatchHdmi(tvStatus)

	devices.DiscoverBroadlinkIRDevice()

	// Check the channel values until termination
	for {
		select {

		case soundStatusUpdate := <-soundOutputStatus:
			soundOuputOn = soundStatusUpdate
			soundOutputChar = OFF_SYMBOL
			if soundOuputOn {
				soundOutputChar = ON_SYMBOL
			}

		case tvStatusUpdate := <-tvStatus:
			tvOn = tvStatusUpdate
			tvStateChar = OFF_SYMBOL
			if tvOn {
				tvStateChar = ON_SYMBOL
			}
		}

		airplayStateChar = OFF_SYMBOL
		sourceSignal := SOURCE_AIRPLAY

		if tvOn {
			// TV viewing
			sourceSignal = SOURCE_TV
			devices.SetBrowserUrl(WEATHER_URL)
			devices.ScreenBrightnessManual(0)
		} else if soundOuputOn {
			// Roon playback (intensive music listening)
			sourceSignal = SOURCE_SOUND_OUTPUT
			devices.SetBrowserUrl(ROON_DISPLAY_URL)
			devices.ScreenBrightnessAuto()
		} else {
			// Airplay (general music / sound)
			airplayStateChar = ON_SYMBOL
			devices.SetBrowserUrl(WEATHER_URL)
			devices.ScreenBrightnessAuto()
		}

		fmt.Printf("Updating source: TV %s\tRoon %s\tAirPlay %s\n", tvStateChar, soundOutputChar, airplayStateChar)
		go devices.SendIRCommand(sourceSignal)
	}
}
