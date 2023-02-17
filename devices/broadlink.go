package devices

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
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

func DiscoverBroadlinkIRDevice() {
	discoverCmd := cmd.NewCmd("broadlink_discovery")
	discoverCmdResult := <-discoverCmd.Start()
	if len(discoverCmdResult.Stderr) > 0 {
		fmt.Println("ERROR: Could not find Broadlink Device.")
	}

	for _, stdOutLine := range discoverCmdResult.Stdout {
		matches := broadlinkDeviceArgsRegex.FindStringSubmatch(stdOutLine)
		if len(matches) == 4 {
			foundDevice = BroadlinkDevice{
				Type: matches[1],
				Host: matches[2],
				Mac:  matches[3],
			}
			broadlinkDevArgs = strings.Fields(fmt.Sprintf(`--type %s --host %s --mac %s`, foundDevice.Type, foundDevice.Host, foundDevice.Mac))
			go TrackTemprature()
		}
	}
}

type TemperatureSample struct {
	SampleTime   int64   `json:"sample_time"`
	DeviceName   string  `json:"device_name"`
	SampledValue float64 `json:"sampled_value"`
}

type TemperatureSamples struct {
	Samples []TemperatureSample `json:"samples"`
}

func TrackTemprature() {
	cadence := time.NewTicker(TEMPERATURE_CADENCE)
	for range cadence.C {
		currentTemp := <-GetTemperature()
		currentTime := time.Now()
		sample := TemperatureSample{
			SampleTime:   currentTime.Unix(),
			DeviceName:   foundDevice.Mac,
			SampledValue: currentTemp,
		}
		samples := TemperatureSamples{
			Samples: []TemperatureSample{sample},
		}
		samplesJson, err := json.Marshal(samples)
		if err != nil {
			fmt.Println("Error marshalling JSON for temperature reading")
		}
		go postTemperature(samplesJson)
	}
}

func postTemperature(samplesJson []byte) chan bool {
	TEMPERATURE_REPORT_URL := os.Getenv("TEMPERATURE_REPORT_URL")
	TEMPERATURE_REPORT_TOKEN := os.Getenv("TEMPERATURE_REPORT_TOKEN")

	responded := make(chan bool)
	request, _ := http.NewRequest("POST", TEMPERATURE_REPORT_URL, bytes.NewBuffer(samplesJson))
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", TEMPERATURE_REPORT_TOKEN))
	request.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	go func() {
		response, err := httpClient.Do(request)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()

		if response.StatusCode != 200 {
			fmt.Printf("Temperature sample failed to upload. Error: %d: %s\n", response.StatusCode, response.Status)
			responded <- false
		} else {
			responded <- true
		}
	}()
	return responded
}

func GetTemperature() chan float64 {
	tempReturn := make(chan float64)
	go func() {
		broadlinkCliTemperatureString := append(broadlinkDevArgs, "--temperature")
		broadlinkCliTemperatureCmd := cmd.NewCmd("broadlink_cli", broadlinkCliTemperatureString...)
		broadlinkCliTemperatureResult := <-broadlinkCliTemperatureCmd.Start()
		if len(broadlinkCliTemperatureResult.Stderr) > 0 {
			fmt.Println("ERROR: Could not read temprature.")
			fmt.Println(broadlinkCliTemperatureResult.Stderr)
		}
		tempString := broadlinkCliTemperatureResult.Stdout[0]
		value, err := strconv.ParseFloat(tempString, 64)
		if err != nil {
			fmt.Println("Could not get a temperature reading.")
		}
		tempReturn <- float64(value)
	}()
	return tempReturn
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
