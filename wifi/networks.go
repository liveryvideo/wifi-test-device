package wifi

import (
	"log"
	"net"
)

type NetworkInterface struct {
	Name            string
	Flags           string
	MTU             int
	HardwareAddress string
	Addresses       []NetworkAddress
}

type NetworkAddress struct {
	Name    string
	Address string
}

func GetNetworkInterfaces() ([]NetworkInterface, error) {
	rawInterfaces, interfaceError := net.Interfaces()

	if interfaceError != nil {
		log.Printf("Could not get network status. %s\n", interfaceError)
		return nil, interfaceError
	}

	networkInterfaces := make([]NetworkInterface, len(rawInterfaces))

	for i, rawInterface := range rawInterfaces {
		addresses, err := rawInterface.Addrs()
		networkAddresses := make([]NetworkAddress, len(addresses))
		for c, address := range addresses {
			networkAddresses[c] = NetworkAddress{
				Name:    address.Network(),
				Address: address.String(),
			}
		}

		if err != nil {
			log.Printf("Failed to parse network interface: %s\n", err)
		}

		networkInterfaces[i] = NetworkInterface{
			Name:            rawInterface.Name,
			Flags:           rawInterface.Flags.String(),
			MTU:             rawInterface.MTU,
			HardwareAddress: rawInterface.HardwareAddr.String(),
			Addresses:       networkAddresses,
		}
	}

	return networkInterfaces, nil
}
