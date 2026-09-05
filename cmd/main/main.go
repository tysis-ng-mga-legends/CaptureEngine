package main

import (
	"fmt"
	"log"
	"time"
	"capture_engine/internal/pkcapture"
	"capture_engine/internal/sender"
)

func main() {

	device := "eth0"
	snaplen := int32(1024)
	promisc := false
	timeout := 1 * time.Second

	
	orchestrator := pkcapture.NewOrchestrator()

	go sender.StartFeatureStream(orchestrator)
	err := pkcapture.PacketCapture(device, snaplen, promisc, int(timeout.Seconds()), orchestrator)
	if err != nil {
		log.Fatalf("Error capturing packets: %v", err)
	}

	fmt.Println("Packet capture completed.")
}