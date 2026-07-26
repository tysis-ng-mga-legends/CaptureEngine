package pkcapture

import (
	"capture_engine/internal/ftextract"
)

type PacketData struct {
    SourceIP  string `json:"source_ip"`
    Timestamp int64 `json:"timestamp"`
    Length    int    `json:"length"`
}
type FlowBatch struct {
	ID      FlowID              `json:"flow_id"`
	Features ftextract.Features  `json:"features"`
}

type FlowOrchestrator struct {
	ActiveSequence  map[FlowID][]PacketData
	OutboundChannel chan FlowBatch
}

// This function creates a new instance of FlowOrchestrator with an initialized ActiveSequence map.
func NewOrchestrator() *FlowOrchestrator {
	return &FlowOrchestrator{
		ActiveSequence:  make(map[FlowID][]PacketData),
		OutboundChannel: make(chan FlowBatch, 100),
	}
}

// This method increments the packet count for a given flowID and stores the packet data in the ActiveSequence map.
func (fo *FlowOrchestrator) IncrementPacketCount(flowID FlowID, packet PacketData) {

	canonicalID := flowID.GetNormalized()
	fo.ActiveSequence[canonicalID] = append(fo.ActiveSequence[canonicalID], PacketData{
		SourceIP:  packet.SourceIP,
		Timestamp: packet.Timestamp,
		Length:    packet.Length,
	})

	if len(fo.ActiveSequence[canonicalID]) == 7 {
	
		packetBatch := fo.ActiveSequence[canonicalID]
		
		lengths := make([]int, len(packetBatch))
		isFwd := make([]bool, len(packetBatch))
		timestamps := make([]int64, len(packetBatch))

		initiatorIP := packetBatch[0].SourceIP
		for i, pkt := range packetBatch {
			lengths[i] = pkt.Length
			isFwd[i] = pkt.SourceIP == initiatorIP
			timestamps[i] = pkt.Timestamp
		}

		features := ftextract.Features{}
		features.ExractFeatures(lengths, isFwd, timestamps)

		completedBatch := FlowBatch{
			ID:      canonicalID,
			Features: features,
		}

		// fmt.Printf("\n\n %v", completedBatch)
		// fo.AnalyzeSequence(canonicalID, packetBatch)
		fo.OutboundChannel <- completedBatch
		
		fo.OutboundChannel <- completedBatch
		// Clear the sequence window so new packets for this flow can be tracked
		delete(fo.ActiveSequence, canonicalID)
	}
}