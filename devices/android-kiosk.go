package devices

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var KIOSK_COMMAND_TEMPLATE = `http://%s/?cmd=%s&password=%s&type=json`

var (
	KIOSK_HOST_AND_PORT = os.Getenv("KIOSK_HOST_AND_PORT")
	KIOSK_PASSWORD      = os.Getenv("KIOSK_PASSWORD")
)

type DeviceInfo struct {
	IsPlugged      bool   `json:"isPlugged"`
	CurrentPageURL string `json:"currentPageUrl"`
	IsInDaydream   bool   `json:"isInDaydream"`
}

func sendCommand(commandUrl string) {
	_, err := http.Get(commandUrl)
	if err != nil {
		fmt.Println(err)
	}
}

func ScreenBrightnessAuto() {
	manualBrightnessCmd := fmt.Sprintf("setStringSetting&key=screenBrightness&value=%s", "")
	setBrightnessUrl := fmt.Sprintf(KIOSK_COMMAND_TEMPLATE, KIOSK_HOST_AND_PORT, manualBrightnessCmd, KIOSK_PASSWORD)
	sendCommand(setBrightnessUrl)
}

func ScreenBrightnessManual(brightness uint) {
	manualBrightnessCmd := fmt.Sprintf("setStringSetting&key=screenBrightness&value=%d", brightness)
	setBrightnessUrl := fmt.Sprintf(KIOSK_COMMAND_TEMPLATE, KIOSK_HOST_AND_PORT, manualBrightnessCmd, KIOSK_PASSWORD)
	sendCommand(setBrightnessUrl)
}

// We want to be fairly protective of this, so that it only gets triggered if the
// user is hands-off the device.

func SetBrowserUrl(navigationUrl string) {
	deviceInfo, err := getDeviceInfo()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("error retrieving device information to gate URL: %v", err))
		return
	}

	// Bail here if the device is not plugged in to charge (i.e. in somebody's hands)
	if !deviceInfo.IsPlugged {
		return
	}

	// Bail here if we're navigating to the same URL
	currentUrl := deviceInfo.CurrentPageURL
	if CompareURLsIgnoringTrailingSlash(currentUrl, navigationUrl) {
		return
	}

	escapedUrl := url.QueryEscape(navigationUrl)
	urlCmd := fmt.Sprintf("loadUrl&url=%s&newtab=false&focus=true", escapedUrl)
	browserNavigateUrl := fmt.Sprintf(KIOSK_COMMAND_TEMPLATE, KIOSK_HOST_AND_PORT, urlCmd, KIOSK_PASSWORD)
	sendCommand(browserNavigateUrl)
}

func getDeviceInfo() (DeviceInfo, error) {
	var deviceInfo DeviceInfo
	deviceInfoUrl := fmt.Sprintf(KIOSK_COMMAND_TEMPLATE, KIOSK_HOST_AND_PORT, "getDeviceInfo", KIOSK_PASSWORD)

	// Send the HTTP GET request
	resp, err := http.Get(deviceInfoUrl)
	if err != nil {
		return deviceInfo, fmt.Errorf("error sending GET request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return deviceInfo, fmt.Errorf("error reading response body: %v", err)
	}

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(body, &deviceInfo)
	if err != nil {
		return deviceInfo, fmt.Errorf("error unmarshalling JSON response: %v", err)
	}

	return deviceInfo, nil
}

// CompareURLsIgnoringTrailingSlash compares two URLs while ignoring trailing slashes
func CompareURLsIgnoringTrailingSlash(url1, url2 string) bool {
	parsedURL1, err1 := url.Parse(url1)
	parsedURL2, err2 := url.Parse(url2)

	if err1 != nil || err2 != nil {
		// Error occurred while parsing URLs, consider them unequal
		return false
	}

	// Remove trailing slashes from the URL paths
	path1 := strings.TrimRight(parsedURL1.Path, "/")
	path2 := strings.TrimRight(parsedURL2.Path, "/")

	// Compare the modified URLs
	return parsedURL1.Scheme == parsedURL2.Scheme &&
		parsedURL1.Host == parsedURL2.Host &&
		path1 == path2
}
