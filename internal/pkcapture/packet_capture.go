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
		srcIP, dstIP, srcPort, dstPort, protocol := CaptureFiveTuples(packet)


		fmt.Printf("Packet captured at %s, SrcIP: %s DstIP: %s SrcPort: %d DstPort: %d Protocol: %s length: %d bytes\n",
		 metadata.Timestamp.Format(time.RFC3339), srcIP, dstIP, srcPort, dstPort, protocol, metadata.Length)	}

	return nil
}

func CaptureFiveTuples(packet gopacket.Packet) (srcIP, dstIP string, srcPort, dstPort uint16, protocol string) {

	netlayer := packet.NetworkLayer()
	if netlayer == nil {
		return "", "", 0, 0, ""
	}

	netflow := netlayer.NetworkFlow()
	srcIP = netflow.Src().String()
	dstIP = netflow.Dst().String()
	
	transportlayer := packet.TransportLayer()
	if transportlayer == nil {
		return srcIP, dstIP, 0, 0, ""
	}

	transFlow := transportlayer.TransportFlow()
	srcPort = uint16(transFlow.Src().EndpointType())
	dstPort = uint16(transFlow.Dst().EndpointType())
	protocol = transportlayer.LayerType().String()

	return srcIP, dstIP, srcPort, dstPort, protocol

}
