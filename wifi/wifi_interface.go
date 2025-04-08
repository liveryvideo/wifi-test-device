package wifi

import (
	"strconv"
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
	if globalRules.LatencyRule.BaseLatency > 0 {
		args = append(args, "delay", baseLatency, variation, correlation)
	}

	// Loss / Corruption / Duplication Rule
	loss := strconv.Itoa(globalRules.Loss) + "%"
	corruption := strconv.Itoa(globalRules.Corruption) + "%"
	duplication := strconv.Itoa(globalRules.Duplication) + "%"
	if globalRules.Loss > 0 {
		args = append(args, "loss", loss)
	}
	if globalRules.Corruption > 0 {
		args = append(args, "corrupt", corruption)
	}
	if globalRules.Duplication > 0 {
		args = append(args, "duplicate", duplication)
	}

	rate := globalRules.BandwidthRule.Rate
	burst := globalRules.BandwidthRule.Burst
	maxLatency := strconv.Itoa(globalRules.BandwidthRule.MaxLatency) + "ms"

	printOutput(performCommand("tc", "qdisc", "add", "dev", "wlan0", "root", "handle", "1:0", "tbf", "rate", rate, "burst", burst, "latency", maxLatency))
	printOutput(performCommand("tc", args...))

}
