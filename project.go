package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var (
	fivePercentAlert   = true
	tenPercentAlert    = true
	twentyPercentAlert = true
	fullChargeAlert    = true
	prevCharging       = -1
	myApp              fyne.App
)

func sendNotification(message string) {
	myApp.SendNotification(&fyne.Notification{
		Title:   "Battery Status",
		Content: message,
	})
}
func getChargingStatus() int {
	out, err := exec.Command("cat", "/sys/class/power_supply/ACAD/online").Output()
	if err != nil {
		fmt.Println("Charging status error:", err)
		return -1
	}
	val, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		fmt.Println("Conversion error (charging status):", err)
		return -1
	}
	return val
}
func getBatteryPercentage() int {
	pathOut, err := exec.Command("bash", "-c", "upower -e | grep BAT").Output()
	if err != nil {
		fmt.Println("Battery path error:", err)
		return -1
	}
	batPath := strings.TrimSpace(string(pathOut))
	if batPath == "" {
		fmt.Println("Battery path not found")
		return -1
	}
	infoOut, err := exec.Command("upower", "-i", batPath).Output()
	if err != nil {
		fmt.Println("Battery percentage error:", err)
		return -1
	}
	for _, line := range strings.Split(string(infoOut), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "percentage:") {
			percentStr := strings.TrimSuffix(
				strings.TrimSpace(strings.TrimPrefix(line, "percentage:")), "%")
			val, err := strconv.Atoi(percentStr)
			if err != nil {
				fmt.Println("Conversion error (battery percentage):", err)
				return -1
			}
			return val
		}
	}
	return -1
}
func batteryPollLoop() {
	for {
		charging := getChargingStatus()
		level := getBatteryPercentage()
		if prevCharging != -1 && charging != prevCharging {
			if charging == 1 {
				fullChargeAlert = true
			} else {
				twentyPercentAlert = true
			}
		}
		prevCharging = charging
		if charging == -1 || level == -1 {
			sendNotification("Error: Unable to read battery info.")
		} else if charging == 0 {
			switch {
			case level <= 5 && fivePercentAlert:
				sendNotification("Battery critical (≤5%)! Please charge now.")
				fivePercentAlert = false
			case level == 10 && tenPercentAlert:
				sendNotification("Battery low (10%) — plug in your charger.")
				tenPercentAlert = false
			case level <= 20 && twentyPercentAlert:
				sendNotification("Battery low (≤20%).")
				twentyPercentAlert = false
			}
		} else {
			if level >= 95 || level == 100 && fullChargeAlert {
				sendNotification(fmt.Sprintf("Battery sufficiently charged %d%%. You can unplug.", level))
				fullChargeAlert = false
			}
		}
		time.Sleep(100 * time.Second)
	}
}
func main() {
	myApp = app.New()
	w := myApp.NewWindow("")
	w.Hide()
	go batteryPollLoop()
	myApp.Run()
}
