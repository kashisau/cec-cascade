package devices

import (
	"fmt"
	"net/http"
	"os"
)

var KIOSK_COMMAND_TEMPLATE = `http://%s:2323/?cmd=%s&password=%s&type=json`

func sendCommand(commandUrl string) {
	_, err := http.Get(commandUrl)
	if err != nil {
		fmt.Println(err)
	}
}

func TurnKioskOn() {
	KIOSK_IP := os.Getenv("KISOK_IP")
	KIOSK_PASSWORD := os.Getenv("KIOSK_PASSWORD")

	screensaverOffUrl := fmt.Sprintf(KIOSK_COMMAND_TEMPLATE, KIOSK_IP, "stopDaydream", KIOSK_PASSWORD)
	sendCommand(screensaverOffUrl)
}

func TurnKioskOff() {
	KIOSK_IP := os.Getenv("KISOK_IP")
	KIOSK_PASSWORD := os.Getenv("KIOSK_PASSWORD")

	screensaverOnUrl := fmt.Sprintf(KIOSK_COMMAND_TEMPLATE, KIOSK_IP, "startDaydream", KIOSK_PASSWORD)
	sendCommand(screensaverOnUrl)
}
