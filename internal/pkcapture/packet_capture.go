package pkcapture


import (
	"fmt"
	"log"
	"time"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)


func PacketCapture(device string, snaplen int32, promisc bool, timeout int) error {

	handle, err := pcap.OpenLive(device, snaplen, promisc, time.Duration(timeout))
	
	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("error opening device %s: %v", device, err)
	}

	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	
	for packet := range packetSource.Packets() {
		
		metadata := packet.Metadata()
		fmt.Printf("Packet captured at %s, length: %d bytes\n", metadata.Timestamp.Format(time.RFC3339), metadata.Length)
}