package wifi

import (
	"log"
	"strconv"
	"strings"
)

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

func GetNetworkStatus() ([]NetworkStatus, error) {
	out, commandError := performCommand("ifconfig")

	if commandError != nil {
		log.Printf("Could not get network status. %s\r\n", commandError)
		return nil, commandError
	}

	rawNetworks := strings.Split(string(out), "\n\n")

	return parseNetworks(rawNetworks), nil
}

func parseNetworks(rawNetworks []string) []NetworkStatus {
	networks := make([]NetworkStatus, len(rawNetworks)-1)
	for i := range networks {
		lines := strings.Split(rawNetworks[i], "\n")
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
