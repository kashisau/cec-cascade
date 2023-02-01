package devices

import (
	"time"

	"github.com/go-cmd/cmd"
)

const CADENCE = 500 * time.Millisecond

var cadence *time.Ticker
var watchSoundCmd *cmd.Cmd

var watchSoundDoneChan = make(chan struct{})
var soundState bool

func WatchSoundOutput(soundStateOut chan<- bool) {
	defer close(watchSoundDoneChan)
	cadence = time.NewTicker(CADENCE)
	for range cadence.C {
		watchSoundCmd = cmd.NewCmd("cat", "/proc/asound/card0/pcm0p/sub0/status")
		watchSoundStatusChan := watchSoundCmd.Start()
		finalStatus := <-watchSoundStatusChan
		newSoundState := len(finalStatus.Stdout) > 1
		if newSoundState != soundState {
			soundStateOut <- newSoundState
			soundState = newSoundState
		}
	}
}
