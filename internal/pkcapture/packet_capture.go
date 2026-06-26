package pkcapture

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)


type FlowID struct {
	SrcIP  string
	DstIP string
	SrcPort uint16
	DstPort uint16
	Protocol string
}

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
		flowID := CaptureFiveTuples(packet)


		fmt.Printf("Packet captured at %s, SrcIP: %s DstIP: %s SrcPort: %d DstPort: %d Protocol: %s length: %d bytes\n",
		metadata.Timestamp.Format(time.RFC3339), flowID.SrcIP, flowID.DstIP, flowID.SrcPort, flowID.DstPort, flowID.Protocol, metadata.Length)	}

	return nil
}

func CaptureFiveTuples(packet gopacket.Packet) (flowID FlowID) {

	if netlayer := packet.NetworkLayer() ; netlayer != nil {
		netflow := netlayer.NetworkFlow()
		flowID.SrcIP = netflow.Src().String()
		flowID.DstIP = netflow.Dst().String()
	}

	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		flowID.SrcPort = uint16(tcp.SrcPort)
		flowID.DstPort = uint16(tcp.DstPort)
		flowID.Protocol = "TCP"
		return flowID
	}

	if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		flowID.SrcPort = uint16(udp.SrcPort)
		flowID.DstPort = uint16(udp.DstPort)
		flowID.Protocol = "UDP"
		return flowID
	}
	
	if transportLayer := packet.TransportLayer(); transportLayer != nil {
		flowID.Protocol = transportLayer.LayerType().String()
	}

	return flowID

}
