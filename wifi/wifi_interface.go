package wifi

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type LatencyRule struct {
	BaseLatency int
	Variation   int
	Correlation int
}

type BandwidthRule struct {
	Rate       string
	Burst      string
	MaxLatency int
}

type GlobalRules struct {
	Loss          int
	Corruption    int
	Duplication   int
	LatencyRule   LatencyRule
	BandwidthRule BandwidthRule
}

var globalRules *GlobalRules

type Device struct {
	ExpirationTime   int
	LinkAddress      string
	HostName         string
	ClientIdentifier string
	Connected        bool
}

type NetworkStatus struct {
	Name       string
	Inet       string
	Netmask    string
	Broadcast  string
	Inet6      string
	PrefixLen  int
	ScopeId    string
	Origin     string
	TxQueueLen int
	OriginType string
}

func performCommand(command string, args ...string) ([]byte, error) {
	fmt.Print("Performing command: " + command + " ")
	fmt.Println(args)
	cmd := exec.Command(command, args...)
	return cmd.CombinedOutput()
}

func checkDevicesConnectionStatus(devices *[]Device) {
	out, err := performCommand("arp")
	if err != nil {
		fmt.Printf("Failed to run command: %s\r\n", err)
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
	out, err := performCommand("cat", "/var/lib/misc/dnsmasq.leases")

	if err != nil {
		fmt.Printf("Failed to run command: %s\r\n", err)
	}

	raw := string(out)
	lines := strings.Split(raw, "\n")

	devices := make([]Device, len(lines)-1)

	for i, line := range lines {
		rawDevice := strings.Fields(line)

		if len(rawDevice) < 4 {
			fmt.Printf("Failed to parse device: %s\r\n", "Invalid number of arguments.")
			continue
		}

		expirationTime, err := strconv.Atoi(rawDevice[0])

		if err != nil {
			fmt.Printf("Failed to parse device: %s\r\n", err)
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

func GetGlobalRules() *GlobalRules {
	if globalRules == nil {
		globalRules = &GlobalRules{
			Loss:        0,
			Corruption:  0,
			Duplication: 0,
			LatencyRule: LatencyRule{
				BaseLatency: 0,
				Variation:   0,
				Correlation: 0,
			},
			BandwidthRule: BandwidthRule{
				Rate:       "1gbit",
				Burst:      "5mbit",
				MaxLatency: 5000,
			},
		}
		UpdateGlobalRules()
	}
	return globalRules
}

func printOutput(output []byte, err error) {
	fmt.Println(string(output))
	if err != nil {
		fmt.Println(err)
	}
}

func UpdateGlobalRules() {
	globalRules := GetGlobalRules()
	// Remove all existing rules for wlan0
	printOutput(performCommand("tc", "qdisc", "del", "dev", "wlan0", "root"))

	// tc class add dev eth1 parent 1: classid 0:1 htb rate 200kbit
	// tc qdisc add dev eth1 parent 1:1 handle 10: netem delay 400000 5 loss 0.03%

	args := []string{"qdisc", "add", "dev", "wlan0", "parent", "1:0", "netem"}

	// Latency Rule
	baseLatency := strconv.Itoa(globalRules.LatencyRule.BaseLatency) + "ms"
	variation := strconv.Itoa(globalRules.LatencyRule.Variation) + "ms"
	correlation := strconv.Itoa(globalRules.LatencyRule.Correlation) + "%"
	if globalRules.LatencyRule.BaseLatency > 10 && globalRules.LatencyRule.Correlation >= 1 {
		args = append(args, "delay", baseLatency, variation, correlation)
	}

	// Loss / Corruption / Duplication Rule
	loss := strconv.Itoa(globalRules.Loss) + "%"
	corruption := strconv.Itoa(globalRules.Corruption) + "%"
	duplication := strconv.Itoa(globalRules.Duplication) + "%"
	if globalRules.Loss >= 1 {
		args = append(args, "loss", loss)
	}
	if globalRules.Corruption >= 1 {
		args = append(args, "corrupt", corruption)
	}
	if globalRules.Duplication >= 1 {
		args = append(args, "duplicate", duplication)
	}

	rate := globalRules.BandwidthRule.Rate
	burst := globalRules.BandwidthRule.Burst
	maxLatency := strconv.Itoa(globalRules.BandwidthRule.MaxLatency) + "ms"

	printOutput(performCommand("tc", "qdisc", "add", "dev", "wlan0", "root", "handle", "1:0", "tbf", "rate", rate, "burst", burst, "latency", maxLatency))
	printOutput(performCommand("tc", args...))

}

func GetNetworkStatus() ([]NetworkStatus, error) {
	out, commandError := performCommand("ifconfig")

	if commandError != nil {
		fmt.Printf("Could not get network status. %s\r\n", commandError)
		return nil, commandError
	}

	rawNetworks := strings.Split(string(out), "\n\n")

	return parseNetworks(rawNetworks), nil
}

func parseNetworks(rawNetworks []string) []NetworkStatus {
	networks := make([]NetworkStatus, len(rawNetworks)-1)
	for i, rawNetwork := range rawNetworks {
		if i >= len(networks) {
			continue
		}
		lines := strings.Split(rawNetwork, "\n")
		networks[i] = NetworkStatus{}
		for l, line := range lines {
			fields := strings.Fields(line)
			parseLine(&networks[i], l, fields)
		}
	}
	return networks
}

func parseLine(network *NetworkStatus, lineIndex int, fields []string) {
	if len(fields) < 1 {
		return
	}
	switch lineIndex {
	case 0:
		network.Name = fields[0]
		return
	case 1:
		broadcast := "NaN"
		if len(fields) >= 6 {
			broadcast = fields[5]
		}

		network.Inet = fields[1]
		network.Netmask = fields[3]
		network.Broadcast = broadcast
		return
	case 2:
		prefixLen, _ := strconv.Atoi(fields[3])

		network.Inet6 = fields[1]
		network.PrefixLen = prefixLen
		network.ScopeId = fields[5]
		return
	case 3:
		length := len(fields)

		txQueuLen, err := strconv.Atoi(fields[length-2])
		if err != nil {
			txQueuLen, err = strconv.Atoi(fields[length-3])
		}
		if err != nil {
			txQueuLen = 1000
		}

		network.Origin = fields[0]
		network.OriginType = fields[length-1]
		network.TxQueueLen = txQueuLen
	}
}
