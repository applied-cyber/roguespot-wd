package main

import (
	"fmt"
	"log"
	"os"
	"time"
	scanner "simple-wd/ap-scanner"
)

func main() {
	// Scanning for access points on Linux needs root. Necessary for other platforms?
	if os.Geteuid() != 0 {
		log.Fatal("The wardriver program needs to be run as root")
	}

	iw, err := scanner.NewIWCommand()
	if err != nil {
		log.Fatal(err)
	}


	// Rescanning over interval
	INTERVAL := 30 * time.Second;
	for range time.Tick(INTERVAL) {
    	
		accessPoints := iw.GetAccessPoints()
	
		// Sanity check: print data
		for _, accessPoint := range accessPoints {
			fmt.Printf(
				"SSID: %s, MAC Address: %s, Signal Strength: %.2f dBm\n",
				accessPoint.SSID,
				accessPoint.Address,
				accessPoint.Strength,
			)
		}
	}

}
