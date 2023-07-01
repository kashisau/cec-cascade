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

type BroadlinkDevice struct {
	Type string
	Host string
	Mac  string
}

var foundDevice BroadlinkDevice
var broadlinkDevArgs []string
var broadlinkDeviceArgsRegex, _ = regexp.Compile(`\#\sbroadlink_cli\s--type\s(0x[0-9a-f]+)\s--host\s([0-9\.]+)\s--mac\s([0-9a-f]+)`)

// DiscoverBroadlinkIRDevice is responsible for discovering a specific Broadlink IR device and populating
// the necessary information about that device. It utilizes an external command-line tool called
// "broadlink_discovery" to perform the device discovery.
func DiscoverBroadlinkIRDevice() {
	TARGET_BROADLINK_MAC := os.Getenv("TARGET_BROADLINK_MAC")
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

	// Intentional panic. The service manager should re-start this executable to reinitialise.
	panic("No Broadlink device to control.")
}

// SendIRCommand sends an IR command using a Broadlink device. It executes the command-line tool "broadlink_cli"
// with the provided sendValue, handling any errors and triggering a panic to signal the need for device re-discovery.
func SendIRCommand(sendValue string) {
	broadlinkCliCmdString := append(broadlinkDevArgs, "--send", sendValue)
	sendCmd := cmd.NewCmd("broadlink_cli", broadlinkCliCmdString...)
	sendCmdResult := <-sendCmd.Start()
	if len(sendCmdResult.Stderr) > 0 {
		fmt.Println("ERROR: Could not send command.")
		fmt.Println(sendCmdResult.Stderr)

		// Intentional panic. The service manager should re-start this executable to reinitialise.
		panic("We need to find the Broadlink device again!")
	}
}
