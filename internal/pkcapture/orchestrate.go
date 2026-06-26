package pkcapture

import (
	"fmt"
	"time"
)

type PacketData struct {
	Timestamp time.Time
	Length int 
}

type FlowOrchestrator struct {
	ActiveSequence map[FlowID][]PacketData
}

func NewOrchestrator() *FlowOrchestrator {
	return &FlowOrchestrator{
		ActiveSequence: make(map[FlowID][]PacketData),
	}
}

func (fo *FlowOrchestrator) IncrementPacketCount(flowID FlowID, packet PacketData) {

	canonicalID := flowID.GetNormalized()

	fo.ActiveSequence[canonicalID] = append(fo.ActiveSequence[canonicalID], PacketData{
		Timestamp: packet.Timestamp,
		Length:    packet.Length,
	})

	if len(fo.ActiveSequence[canonicalID]) == 5 {
	
		packetBatch := fo.ActiveSequence[canonicalID]

		fo.AnalyzeSequence(canonicalID, packetBatch)
		
		// delete (fo.ActiveSequence, canonicalID)
	}
}

func (o *FlowOrchestrator) AnalyzeSequence(id FlowID, packets []PacketData) {
	fmt.Printf("\n==================================================\n")
	fmt.Printf("FLOW BATCH READY FOR CLASSIFICATION\n")
	fmt.Printf("Protocol: %s | Host A: %s:%d | Host B: %s:%d\n", id.Protocol, id.SrcIP, id.SrcPort, id.DstIP, id.DstPort)
	fmt.Printf("Total Packets in Batch: %d\n", len(packets))

	for idx, packet := range packets {
		fmt.Printf("  Packet %d -> Captured: %s | Size: %d bytes\n", 
			idx+1, 
			packet.Timestamp.Format("15:04:05.000000"), 
			packet.Length,
		)
	}
	fmt.Printf("==================================================\n")
}