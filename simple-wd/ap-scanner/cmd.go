// Runs the 'iw' command on Linux and parses the output to create a list of access points.
// In the future, it is planned to expand to use 'airport' for MacOS and 'net.sh' for Windows.
package scanner

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"strconv"
	"strings"
)

const (
	interfaceName string = "wlan0" // Should eventually move to a configuration file
	iwCommandName string = "iw"
)

type IWCommand struct {
	interfaceName string
}

// isCommandAvailable checks if a command is available in the path
func isCommandAvailable(command string) bool {
	if _, err := exec.LookPath(command); err != nil {
		return false
	}
	return true
}

// runCommand runs a command with the command and its arguments. Returns
// an error if the command is not found. Exits if the command returns
// an exit code >= 1
func runCommand(command string, args ...string) (string, error) {
	if !isCommandAvailable(command) {
		return "", &CommandNotFoundError{command}
	}

	cmd := exec.Command(command, args...)
	cmdOutput, err := cmd.CombinedOutput()
	output := string(cmdOutput)

	if err != nil {
		// Log the error and exit
		log.Fatalf(
			"Command '%s %s' failed to run: '%s'",
			command,
			strings.Join(args, " "),
			strings.TrimSuffix(output, "\n"),
		)
	}

	return output, nil
}

// isValidInterfaceName checks if the interface name provided is a valid interface
// that is detectable by Go
func isValidInterfaceName(name string) (bool, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return false, err
	}

	for _, iface := range ifaces {
		if iface.Name == name {
			return true, nil
		}
	}

	return false, nil
}

// NewIWCommand does a few checks to ensure the command can run successfully before
// creating the command
func NewIWCommand() (*IWCommand, error) {
	if !isCommandAvailable(iwCommandName) {
		return nil, &CommandNotFoundError{iwCommandName}
	}

	isValid, err := isValidInterfaceName(interfaceName)
	if err != nil {
		return nil, err
	}

	if !isValid {
		return nil, fmt.Errorf("Interface '%s' not found", interfaceName)
	}

	iw := &IWCommand{interfaceName}
	return iw, nil
}

// scan scans for access points
func (iw *IWCommand) scan() (string, error) {
	cmdOutput, err := runCommand("iw", "dev", iw.interfaceName, "scan")
	if err != nil {
		return "", err
	}

	return cmdOutput, nil
}

// parseScan parses the output from the scan command and returns a list
// of access points, each with an MAC address, SSID, and signal strength
func (iw *IWCommand) parseScan(output string) []APInfo {
	accessPoints := make([]APInfo, 0)

	macAddress, ssid, signalStrength := "", "", 0.0
	createAP := false

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		// Trim leading and trailing whitespace from the line
		trimmedLine := strings.TrimSpace(line)
		fields := strings.Fields(line)
		if len(fields) < 2 {
			// The line doesn't have enough info, so don't attempt to parse it
			continue
		}

		if strings.HasPrefix(line, "BSS") {
			// MAC address
			macAddress = strings.TrimSuffix(fields[1], "(on")
		} else if strings.HasPrefix(trimmedLine, "SSID:") {
			// SSID
			ssid = fields[1]
			// SSID is the last of the three access point attributes we want, so create an
			// AP after this block
			createAP = true
		} else if strings.HasPrefix(trimmedLine, "signal:") {
			// Signal strength
			strength := strings.TrimSuffix(fields[1], "dBm")
			signalStrength, _ = strconv.ParseFloat(strength, 64)
		}

		if createAP {
			if macAddress != "" && ssid != "" && signalStrength != 0 {
				// Only add an access point if all attributes were set
				accessPoint := APInfo{macAddress, ssid, signalStrength}
				accessPoints = append(accessPoints, accessPoint)
			}

			// Reset the values for the next access point
			macAddress, ssid, signalStrength = "", "", 0.0
			createAP = false
		}
	}

	return accessPoints
}

// GetAccessPoints scans for access points and parses them into an easy-to-handle
// list of access points. Returns an empty list of APs if the scan fails
func (iw *IWCommand) GetAccessPoints() []APInfo {
	scanOutput, err := iw.scan()
	if err != nil {
		// In case of an error, skip access point retrieval
		return []APInfo{}
	}

	accessPoints := iw.parseScan(scanOutput)
	return accessPoints
}
