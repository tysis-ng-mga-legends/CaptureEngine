package main

import (
	"capture_engine/internal/pkcapture"
	"capture_engine/internal/sender"
	"fmt"
	"log"
	"time"
)

func main() {
	// Sniff directly on the Wi-Fi hotspot interface
	device := "wlan0"
	snaplen := int32(1024)
	promisc := false
	timeout := 1 * time.Second

	orchestrator := pkcapture.NewOrchestrator()
	go sender.StartFeatureStream(orchestrator)

	fmt.Printf("[*] Starting packet capture on %s...\n", device)
	err := pkcapture.PacketCapture(device, snaplen, promisc, timeout, orchestrator)
	if err != nil {
		log.Fatalf("Error capturing packets: %v", err)
	}
}