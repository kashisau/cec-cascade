package devices

import (
	"fmt"
	"os"
	"regexp"

	"github.com/go-cmd/cmd"
)

var cecCmd *cmd.Cmd
var cecDoneChan = make(chan struct{})
var tvState bool
var tvStateRegex *regexp.Regexp

func WatchHdmi(tvStateChangeOut chan<- bool) {
	tvStateRegex, _ = regexp.Compile(`TV\s\(\d\)\:\spower status changed from '([a-z\s]+)' to '([a-z\s]+)'`)

	// Disable output buffering, enable streaming
	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}
	cecCmd = cmd.NewCmdOptions(cmdOptions, "cec-client")
	cecCmd.Start()
	defer close(cecDoneChan)

	for cecCmd.Stdout != nil || cecCmd.Stderr != nil {
		select {
		case stdOutLine := <-cecCmd.Stdout:
			matches := tvStateRegex.FindStringSubmatch(stdOutLine)
			if len(matches) == 3 {
				_, to := matches[1], matches[2]
				newTVState := to == "on"
				if newTVState != tvState {
					tvStateChangeOut <- newTVState
					tvState = newTVState
				}
			}
		case line, open := <-cecCmd.Stderr:
			if !open {
				cecCmd.Stderr = nil
				continue
			}
			fmt.Fprintln(os.Stderr, line)
		}
	}
}
