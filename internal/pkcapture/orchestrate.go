package pkcapture

import (
	"capture_engine/internal/ftextract"
	"fmt"
)

type PacketData struct {
    SourceIP  string `json:"source_ip"`
	Protocol string `json:"proto"`
    Timestamp int64 `json:"timestamp"`
    Length    int    `json:"length"`
	HasPSH    bool   `json:"has_psh"`
}
type FlowBatch struct {
	ID      FlowID              `json:"flow_id"`
	Protocol string 			`json:"protocol"`
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
	fo.ActiveSequence[canonicalID] = append(fo.ActiveSequence[canonicalID], packet)

	if len(fo.ActiveSequence[canonicalID]) == 7 {
	
		packetBatch := fo.ActiveSequence[canonicalID]

		resolvedProto := canonicalID.Protocol 
		for _, p := range packetBatch{
			if p.Protocol == "TLS" || p.Protocol == "QUIC" {
				resolvedProto = p.Protocol
				break
			}
			fmt.Println(p.Protocol)
		}

		lengths := make([]int, len(packetBatch))
		isFwd := make([]bool, len(packetBatch))
		timestamps := make([]int64, len(packetBatch))
		pshFlags := make([]bool, len(packetBatch))

		initiatorIP := packetBatch[0].SourceIP
		for i, pkt := range packetBatch {
			lengths[i] = pkt.Length
			isFwd[i] = pkt.SourceIP == initiatorIP
			timestamps[i] = pkt.Timestamp
			pshFlags[i] = pkt.HasPSH
		}


		features := ftextract.ExtractFeatures(lengths, isFwd, timestamps, pshFlags, resolvedProto)

		completedBatch := FlowBatch{
			ID:      canonicalID,
			Protocol: resolvedProto,
			Features: features,
		}
		
		fo.OutboundChannel <- completedBatch
		// Clear the sequence window so new packets for this flow can be tracked
		// delete(fo.ActiveSequence, canonicalID)
	}
}