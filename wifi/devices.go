package wifi

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type Device struct {
	ExpirationTime   int
	LinkAddress      string
	HostName         string
	ClientIdentifier string
	Connected        bool
}

func checkDevicesConnectionStatus(devices *[]Device) {
	out, err := performCommand("arp")
	if err != nil {
		log.Printf("Failed to run command: %s\r\n", err)
	}

	raw := string(out)
	lines := strings.Split(raw, "\n")

	for _, line := range lines {
		arpFields := strings.Fields(line)
		if len(arpFields) < 3 {
			continue
		}
		for deviceIndex, _ := range *devices {
			if arpFields[2] == (*devices)[deviceIndex].LinkAddress {
				(*devices)[deviceIndex].Connected = true
				break
			}
		}
	}
}

func GetLeasedDevices() []Device {
	file, err := os.Open("/var/lib/misc/dnsmasq.leases")

	if err != nil {
		log.Printf("Failed to open file: %s\r\n", err)
	}

	out, err := ioutil.ReadAll(file)

	if err != nil {
		log.Printf("Failed to read file: %s\r\n", err)
	}

	raw := string(out)
	lines := strings.Split(raw, "\n")

	devices := make([]Device, len(lines)-1)

	for i := range devices {
		rawDevice := strings.Fields(lines[i])

		if len(rawDevice) < 4 {
			log.Printf("Failed to parse device: %s\r\n", "Invalid number of arguments.")
			continue
		}

		expirationTime, err := strconv.Atoi(rawDevice[0])

		if err != nil {
			log.Printf("Failed to parse device: %s\r\n", err)
			continue
		}

		devices[i] = Device{
			ExpirationTime:   expirationTime,
			LinkAddress:      rawDevice[1],
			HostName:         rawDevice[2],
			ClientIdentifier: rawDevice[3],
			Connected:        false,
		}
	}
	checkDevicesConnectionStatus(&devices)
	return devices
}
