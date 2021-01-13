package main

import (
	"wifi-test-device/rest"
	"wifi-test-device/wifi"
)

func main() {
	wifi.UpdateGlobalRules()
	rest.StartServer(80)
}
