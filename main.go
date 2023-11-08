package main

import (
    "strings"
    "fmt"
    "strconv"
    "io/ioutil"
)

//TODO:
// func runIW(interfaceName string) (string, error) {
//     cmd := exec.Command("iw", "dev", interfaceName, "scan")
//     output, err := cmd.CombinedOutput()
//     if err != nil {
//         return "", err
//     }
//     return string(output), nil
// }

func parseIWOutput(iwOutput string) ([]string, []string, []float64) {
    var macAddresses []string
    var ssids []string
    var signalStrengths []float64

    // Split iwOutput into lines.
    lines := strings.Split(iwOutput, "\n")
    for _, line := range lines {
        // Trim leading and trailing whitespace from the line.
        trimmedLine := strings.TrimSpace(line)
      
        if strings.HasPrefix(line, "BSS") {
            fields := strings.Fields(line)
            if len(fields) >= 2 {
                mac := strings.TrimSuffix(fields[1], "(on")
                macAddresses = append(macAddresses, mac)
            }
        } else if strings.HasPrefix(trimmedLine, "SSID:") {
            fields := strings.Fields(line)
            if len(fields) >= 2 {
                ssids = append(ssids, fields[1])
            }
        } else if strings.HasPrefix(trimmedLine, "signal:") {
            fields := strings.Fields(line)
            if len(fields) >= 2 {
                signalStrength := strings.TrimSuffix(fields[1], "dBm")
                strength, _ := strconv.ParseFloat(signalStrength, 64)
                signalStrengths = append(signalStrengths, strength)
            }
        }
    }
    return macAddresses, ssids, signalStrengths
}

func main() {
    // interfaceName := "wlan0"
    // iwOutput, err := runIW(interfaceName)
    // if err != nil {
    //     panic(err)
    // }

    //I directly open iw.output to parse, should be replaced by the above
    //when runIW is implemented
    filePath := "iw.output"
    fileContent, _ := ioutil.ReadFile(filePath)
    iwOutput := string(fileContent)

    //data to be processed
    macAddresses, ssids, signalStrengths := parseIWOutput(iwOutput)

    //sanity check: print data
    for i := 0; i < len(macAddresses); i++ {
        mac := macAddresses[i]
        ssid := ssids[i]
        signalStrength := signalStrengths[i]
        fmt.Printf("SSID: %s, MAC Address: %s, Signal Strength: %.2f dBm\n", ssid, mac, signalStrength)
    }  
}
