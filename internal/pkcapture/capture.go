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

		orchestrator.IncrementPacketCount(flowID, PacketData{Timestamp: uint64(metadata.Timestamp.UnixNano()), Length: metadata.Length})

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

	hasTCP := false
	hasUDP := false

	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		flowID.SrcPort = uint16(tcp.SrcPort)
		flowID.DstPort = uint16(tcp.DstPort)
		flowID.Protocol = "TCP"
		hasTCP = true
	} else if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		flowID.SrcPort = uint16(udp.SrcPort)
		flowID.DstPort = uint16(udp.DstPort)
		flowID.Protocol = "UDP"
		hasUDP = true
	}
	
	if hasTCP {
		// If it has TCP, check if there is a TLS payload inside it
		if tlsLayer := packet.Layer(layers.LayerTypeTLS); tlsLayer != nil {
			flowID.Protocol = "TLS"
		}
	} else if hasUDP {
		// If it has UDP, check if it is QUIC
		// Note: QUIC packets are encrypted from the first packet, so standard gopacket 
		// doesn't have a default "LayerTypeQUIC" parser. 
		// Industry standard check: UDP traffic on Port 443 is almost always QUIC.
		if flowID.SrcPort == 443 || flowID.DstPort == 443 {
			flowID.Protocol = "QUIC"
		}
	}

	return flowID
}

