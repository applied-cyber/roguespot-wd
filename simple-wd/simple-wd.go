package main

import (
	"fmt"
	"log"
	"os"
	"time"

	scanner "simple-wd/ap-scanner"
	sender "simple-wd/ap-sender"
)

func main() {
	// Scanning for access points on Linux needs root. Necessary for other platforms?
	if os.Geteuid() != 0 {
		log.Fatal("The wardriver program needs to be run as root")
	}

	config := NewConfig()
	iw, err := scanner.NewIWCommand(config.InterfaceName)
	if err != nil {
		log.Fatal(err)
	}

	endpointURL := fmt.Sprintf("http://%s:%d/log", config.Host, config.Port)
	sender := sender.NewSender(endpointURL)

	// Rescanning over interval
	interval := time.Second * time.Duration(config.ScanIntervalSeconds)
	ticker := time.NewTicker(interval)
	for ; true; <-ticker.C {
		accessPoints := iw.GetAccessPoints()

		// Sanity check: print data
		for _, accessPoint := range accessPoints {
			// TODO: Handle errors from send (except failed request errors as they will be enqueued again)
			_ = sender.Send(accessPoint)
			fmt.Printf(
				"SSID: %s, MAC Address: %s, Signal Strength: %.2f dBm\n",
				accessPoint.SSID,
				accessPoint.Address,
				accessPoint.Strength,
			)
		}
	}

}
