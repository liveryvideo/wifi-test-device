package main

import (
	rest "bitbucket.org/exmachina/wifi-test-device/rest"
	wifi "bitbucket.org/exmachina/wifi-test-device/wifi"
)

func main() {
	wifi.UpdateGlobalRules()
	rest.StartServer(80)
}
