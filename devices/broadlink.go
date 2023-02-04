package devices

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-cmd/cmd"
)

var broadlinkDevArgs []string
var broadlinkDeviceArgsRegex, _ = regexp.Compile(`\#\sbroadlink_cli\s(--type\s0x[0-9a-f]+\s--host\s[0-9\.]+\s--mac\s[0-9a-f]+)`)

func DiscoverBroadlinkIRDevice() {
	discoverCmd := cmd.NewCmd("broadlink_discovery")
	discoverCmdResult := <-discoverCmd.Start()
	if len(discoverCmdResult.Stderr) > 0 {
		fmt.Println("ERROR: Could not find Broadlink Device.")
	}

	for _, stdOutLine := range discoverCmdResult.Stdout {
		matches := broadlinkDeviceArgsRegex.FindStringSubmatch(stdOutLine)
		if len(matches) == 2 {
			broadlinkDevArgs = strings.Fields(matches[1])
		}
	}
}

func SendIRCommand(sendValue string) {
	broadlinkCliCmdString := append(broadlinkDevArgs, "--send", sendValue)
	sendCmd := cmd.NewCmd("broadlink_cli", broadlinkCliCmdString...)
	sendCmdResult := <-sendCmd.Start()
	if len(sendCmdResult.Stderr) > 0 {
		fmt.Println("ERROR: Could not send command.")
		fmt.Println(sendCmdResult.Stderr)
	}
}
