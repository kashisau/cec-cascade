package devices

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-cmd/cmd"
)

const TEMPERATURE_CADENCE = 5 * time.Second

var TARGET_BROADLINK_MAC = os.Getenv("TARGET_BROADLINK_MAC")

type BroadlinkDevice struct {
	Type string
	Host string
	Mac  string
}

var foundDevice BroadlinkDevice
var broadlinkDevArgs []string
var broadlinkDeviceArgsRegex, _ = regexp.Compile(`\#\sbroadlink_cli\s--type\s(0x[0-9a-f]+)\s--host\s([0-9\.]+)\s--mac\s([0-9a-f]+)`)

func DiscoverBroadlinkIRDevice() {
	discoverCmd := cmd.NewCmd("broadlink_discovery")
	discoverCmdResult := <-discoverCmd.Start()
	if len(discoverCmdResult.Stderr) > 0 {
		fmt.Println("ERROR: Could not find Broadlink Device.")
	}

	for _, stdOutLine := range discoverCmdResult.Stdout {
		matches := broadlinkDeviceArgsRegex.FindStringSubmatch(stdOutLine)
		if len(matches) == 4 {
			if matches[3] != TARGET_BROADLINK_MAC {
				// Skip this device
				continue
			}

			foundDevice = BroadlinkDevice{
				Type: matches[1],
				Host: matches[2],
				Mac:  matches[3],
			}
			broadlinkDevArgs = strings.Fields(fmt.Sprintf(`--type %s --host %s --mac %s`, foundDevice.Type, foundDevice.Host, foundDevice.Mac))
			return
		}
	}

	fmt.Printf("Could not find the target Broadlink device with machine ID '%s'.", TARGET_BROADLINK_MAC)
	panic("No Broadlink device to control.")
}

func SendIRCommand(sendValue string) {
	broadlinkCliCmdString := append(broadlinkDevArgs, "--send", sendValue)
	sendCmd := cmd.NewCmd("broadlink_cli", broadlinkCliCmdString...)
	sendCmdResult := <-sendCmd.Start()
	if len(sendCmdResult.Stderr) > 0 {
		fmt.Println("ERROR: Could not send command.")
		fmt.Println(sendCmdResult.Stderr)
		panic("We need to find the Broadlink device again!")
	}
}
