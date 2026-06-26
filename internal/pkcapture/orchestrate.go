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

		fmt.Printf("FlowID: %+v, Packet Count: %d, Packets: %+v\n",
			canonicalID, len(packetBatch), packetBatch)
		
		delete (fo.ActiveSequence, canonicalID)
	}
}