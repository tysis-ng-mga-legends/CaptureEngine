package pkcapture

import (
	"capture_engine/internal/ftextract"
	"capture_engine/internal/models"
)

type FlowBatch struct {
	ID      FlowID              `json:"flow_id"`
	Packets []models.PacketData `json:"packets"`
	Features ftextract.SpatialFeatures  `json:"features"`
}

type FlowOrchestrator struct {
	ActiveSequence  map[FlowID][]models.PacketData
	OutboundChannel chan FlowBatch
}

// This function creates a new instance of FlowOrchestrator with an initialized ActiveSequence map.
func NewOrchestrator() *FlowOrchestrator {
	return &FlowOrchestrator{
		ActiveSequence:  make(map[FlowID][]models.PacketData),
		OutboundChannel: make(chan FlowBatch, 100),
	}
}

// This method increments the packet count for a given flowID and stores the packet data in the ActiveSequence map.
func (fo *FlowOrchestrator) IncrementPacketCount(flowID FlowID, packet models.PacketData) {

	canonicalID := flowID.GetNormalized()

	fo.ActiveSequence[canonicalID] = append(fo.ActiveSequence[canonicalID], models.PacketData{
		SourceIP:  packet.SourceIP,
		Timestamp: packet.Timestamp,
		Length:    packet.Length,
	})

	if len(fo.ActiveSequence[canonicalID]) == 7 {
	
		packetBatch := fo.ActiveSequence[canonicalID]
		
		// Call the spatial feature extraction function correctly
		spatialFeats := ftextract.ExtractSpatialFeatures(packetBatch)
		
		completedBatch := FlowBatch{
			ID:      canonicalID,
			Packets: packetBatch,
			Features: spatialFeats,
		}
		
		fo.OutboundChannel <- completedBatch
		
		// Clear the sequence window so new packets for this flow can be tracked
		delete(fo.ActiveSequence, canonicalID)
	}
}