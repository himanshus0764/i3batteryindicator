package main

import (
	"os/exec"
	"strconv"
	"strings"
)

func showZenity(message string) {
	cmd := exec.Command("bash", "-c", `zenity --warning --text="`+message+`"`)
	_ = cmd.Run() // Ignore error if Zenity is closed/cancelled
}
func main() {
	getBat := `upower -e | grep BAT`
	cmdBat, err := exec.Command("bash", "-c", getBat).Output()
	if err != nil {
		showZenity("Error finding battery device: " + err.Error())
		return
	}
	batPath := strings.TrimSpace(string(cmdBat))

	// Query only percentage
	query := `upower -i "` + batPath + `" | grep percentage`
	outBytes, err := exec.Command("bash", "-c", query).Output()
	if err != nil {
		showZenity("Error querying battery info: " + err.Error())
		return
	}

	line := strings.TrimSpace(string(outBytes))
	pctStr := ""
	if strings.HasPrefix(line, "percentage:") {
		pctStr = strings.TrimSuffix(strings.Fields(line)[1], "%")
	}

	// Convert percentage to int
	pct, err := strconv.Atoi(pctStr)
	if err != nil {
		showZenity("Error parsing battery percentage: " + err.Error())
		return
	}

	// Zenity alerts based on percentage
	if pct >= 95 && pct <= 100 {
		showZenity("🔋 Battery is charged – consider unplugging.")
	}

	switch pct {
	case 90:
		showZenity("Battery is almost full (90%)")
	case 20:
		showZenity("Battery is at 20% – please plug in")
	case 10:
		showZenity("Immediately plug in – battery is below 10%")
	}
}
