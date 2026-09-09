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
func PacketCapture(device string, snaplen int32, promisc bool, timeout time.Duration, orchestrator *FlowOrchestrator) error {

	handle, err := pcap.OpenLive(device, snaplen, promisc, timeout) 
	
	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("error opening device %s: %v", device, err)
	}

	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	
	for packet := range packetSource.Packets() {
		
		metadata := packet.Metadata()
		flowID, ProtoType, hasPSH, isAppData, payloadLen := ClassifyAndFilterPacket(packet)

		if !isAppData {
			continue
		}

		packetData :=  PacketData{
			SourceIP: flowID.SrcIP,
			SourcePort: flowID.SrcPort,
			Protocol: string(ProtoType),
			Timestamp: int64(metadata.Timestamp.UnixMicro()),
			Length: metadata.Length,
			PayloadLen: payloadLen,
			HasPSH: hasPSH,
		}
		orchestrator.IncrementPacketCount(flowID, packetData)

	}

	return nil
}

// This function extracts the 5-tuple (source IP, destination IP, source port, destination port, protocol) from a given packet.
func ClassifyAndFilterPacket(packet gopacket.Packet) (FlowID, ProtocolType, bool, bool, int) {
	var flowID FlowID
	netlayer := packet.NetworkLayer()
	if netlayer == nil {
		return flowID, "", false, false, 0
	}
	netflow := netlayer.NetworkFlow()
	flowID.SrcIP = netflow.Src().String()
	flowID.DstIP = netflow.Dst().String()


	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		flowID.SrcPort = uint16(tcp.SrcPort)
		flowID.DstPort = uint16(tcp.DstPort)
		flowID.Protocol = "TCP"

		payload := tcp.Payload
		payloadLen := len(payload)
		hasPSH := tcp.PSH

		if isTLS(payload){
			if isTLSAppData(payload) {
				return flowID, ProtoTLS, hasPSH, true, payloadLen
			} 
			return flowID, ProtoTLS, hasPSH, false, payloadLen
		}
		
		if isPlainTCPAppData(tcp) {
			return flowID, ProtoTCP, hasPSH, true, payloadLen
		}

		return flowID, ProtoTCP, hasPSH, false, payloadLen

	} else if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		flowID.SrcPort = uint16(udp.SrcPort)
		flowID.DstPort = uint16(udp.DstPort)
		flowID.Protocol = "UDP"

		payload := udp.Payload
		payloadLen := len(payload)

		if isQUIC(payload) {
			if isQUICAppData(payload) {
				return flowID, ProtoQUIC, false, true, payloadLen
			}
			return flowID, ProtoQUIC, false, false, payloadLen
		}
		
		if len(payload) > 0 {
			return flowID, ProtoUDP, false, true, payloadLen
		}

		return flowID, ProtoUDP, false, false, payloadLen
	}
	
	return flowID, "", false, false, 0
}

