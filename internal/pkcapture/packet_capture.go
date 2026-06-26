package pkcapture

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
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

	if netlayer := packet.NetworkLayer() ; netlayer != nil {
		netflow := netlayer.NetworkFlow()
		srcIP = netflow.Src().String()
		dstIP = netflow.Dst().String()
	}

	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		srcPort := uint16(tcp.SrcPort)
		dstPort := uint16(tcp.DstPort)
		protocol = "TCP"
		return srcIP, dstIP, srcPort, dstPort, protocol
	}

	if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		srcPort := uint16(udp.SrcPort)
		dstPort := uint16(udp.DstPort)
		protocol = "UDP"
		return srcIP, dstIP, srcPort, dstPort, protocol
	}
	
	if transportLayer := packet.TransportLayer(); transportLayer != nil {
		protocol = transportLayer.LayerType().String()
	}

	return srcIP, dstIP, 0, 0, protocol

}
