package pkcapture

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// This function captures packets from the specified network device and processes them using the FlowOrchestrator.
func PacketCapture(device string, snaplen int32, promisc bool, timeout int, orchestrator *FlowOrchestrator) error {

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

		if flowID.Protocol == "" {
			continue
		}
		
		var hasPSH bool 
		if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
			tcp, _ := tcpLayer.(*layers.TCP)
			hasPSH = (tcp.PSH)
		}

		orchestrator.IncrementPacketCount(flowID, PacketData{SourceIP: flowID.SrcIP, Timestamp: int64(metadata.Timestamp.UnixMicro()), Length: metadata.Length, HasPSH: hasPSH})

	}

	return nil
}

// This function extracts the 5-tuple (source IP, destination IP, source port, destination port, protocol) from a given packet.
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

		if appLayer :=  packet.ApplicationLayer(); appLayer != nil {
			if isTLS(appLayer.Payload()){
				flowID.Protocol = "TLS"
			}
		}
	} else if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		flowID.SrcPort = uint16(udp.SrcPort)
		flowID.DstPort = uint16(udp.DstPort)
		flowID.Protocol = "UDP"

		if appLayer := packet.ApplicationLayer(); appLayer != nil {
			if isQUIC(appLayer.Payload()) {
				flowID.Protocol = "QUIC"
			}
		}
	}
	
	return flowID
}

