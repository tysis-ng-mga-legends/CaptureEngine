package pkcapture

import (
	"capture_engine/internal/ftextract"
)

type PacketData struct {
	Timestamp int64 `json:"timestamp"`
	Length int `json:"length"` 
}

type FlowBatch struct {
	ID FlowID `json:"flow_id"`
	// Packets []PacketData `json:"packets"`
	Features ftextract.Features

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

	if len(fo.ActiveSequence[canonicalID]) == 7 {
	
		packetBatch := fo.ActiveSequence[canonicalID]
		
		timestamp := make([]int64, len(packetBatch))	
		length := make([]int, len(packetBatch))

		for i, p:= range packetBatch {
			timestamp[i] = p.Timestamp
			length[i] = p.Length
		}
		
		ft :=  ftextract.Features{}
		features := ft.ExtractFeatures(timestamp, length)
		
		completedBatch := FlowBatch{
			ID: canonicalID,
			Features: features,
		}

		// fmt.Printf("\n\n %v", completedBatch)
		// fo.AnalyzeSequence(canonicalID, packetBatch)
		fo.OutboundChannel <- completedBatch
		
		// delete (fo.ActiveSequence, canonicalID)
	}
}
