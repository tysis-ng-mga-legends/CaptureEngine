package pkcapture

import (
	"fmt"
	"time"
)

type PacketData struct {
	Timestamp time.Time `json:"timestamp"`
	Length int `json:"length"` 
}

type FlowBatch struct {
	ID FlowID `json:"flow_id"`
	Packets []PacketData `json:"packets"`

}
type FlowOrchestrator struct {
	ActiveSequence map[FlowID][]PacketData
	OutboundChannel chan FlowBatch
}

// This function creates a new instance of FlowOrchestrator with an initialized ActiveSequence map.
func NewOrchestrator() *FlowOrchestrator {
	return &FlowOrchestrator{
		ActiveSequence: make(map[FlowID][]PacketData),
		OutboundChannel: make(chan FlowBatch, 100),
	}
}

// This method increments the packet count for a given flowID and stores the packet data in the ActiveSequence map.
func (fo *FlowOrchestrator) IncrementPacketCount(flowID FlowID, packet PacketData) {

	canonicalID := flowID.GetNormalized()

	fo.ActiveSequence[canonicalID] = append(fo.ActiveSequence[canonicalID], PacketData{
		Timestamp: packet.Timestamp,
		Length:    packet.Length,
	})

	if len(fo.ActiveSequence[canonicalID]) == 15 {
	
		packetBatch := fo.ActiveSequence[canonicalID]

		completedBatch := FlowBatch{
			ID: canonicalID,
			Packets: packetBatch,
		}

		// fo.AnalyzeSequence(canonicalID, packetBatch)
		fo.OutboundChannel <- completedBatch
		
		// delete (fo.ActiveSequence, canonicalID)
	}
}

// This method is just a placheholder, this will be replaced by feature extraction logic.
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