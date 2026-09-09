package pkcapture

import (
	"capture_engine/internal/ftextract"
)

type PacketData struct {
	SourceIP   string `json:"source_ip"`
	SourcePort uint16 `json:"source_port"`
	Protocol   string `json:"proto"`
	Timestamp  int64  `json:"timestamp"`
	Length     int    `json:"length"`
	PayloadLen int    `json:"payload_len"`
	HasPSH     bool   `json:"has_psh"`
}

type FlowBatch struct {
	ID       FlowID             `json:"flow_id"`
	Protocol string             `json:"protocol"`
	Features ftextract.Features `json:"features"`
}

type FlowOrchestrator struct {
	ActiveSequence  map[FlowID][]PacketData
	OutboundChannel chan FlowBatch
}

func NewOrchestrator() *FlowOrchestrator {
	return &FlowOrchestrator{
		ActiveSequence:  make(map[FlowID][]PacketData),
		OutboundChannel: make(chan FlowBatch, 100),
	}
}

func (fo *FlowOrchestrator) IncrementPacketCount(flowID FlowID, packet PacketData) {
	canonicalID := flowID.GetNormalized()

	// Guard against unbounded slice growth (memory leak)
	if len(fo.ActiveSequence[canonicalID]) >= 7 {
		return
	}

	fo.ActiveSequence[canonicalID] = append(fo.ActiveSequence[canonicalID], packet)

	if len(fo.ActiveSequence[canonicalID]) == 7 {
		packetBatch := fo.ActiveSequence[canonicalID]

		// Resolve protocol
		resolvedProto := canonicalID.Protocol
		for _, p := range packetBatch {
			if p.Protocol == "TLS" || p.Protocol == "QUIC" {
				resolvedProto = p.Protocol
				break
			}
		}

		// Single loop to extract raw arrays
		frameLens := make([]float32, 7)
		payloadLens := make([]float32, 7)
		isFwd := make([]bool, 7)
		timestamps := make([]int64, 7)
		pshFlags := make([]bool, 7)

		initiatorIP := packetBatch[0].SourceIP
		initiatorPort := packetBatch[0].SourcePort

		for i, pkt := range packetBatch {
			frameLens[i] = float32(pkt.Length)
			payloadLens[i] = float32(pkt.PayloadLen)
			isFwd[i] = (pkt.SourceIP == initiatorIP && pkt.SourcePort == initiatorPort)
			timestamps[i] = pkt.Timestamp
			pshFlags[i] = pkt.HasPSH
		}

		// Pass frameLens into the compiler
		features := ftextract.ExtractFeatures(frameLens, payloadLens, isFwd, timestamps, pshFlags, resolvedProto)

		completedBatch := FlowBatch{
			ID:       canonicalID,
			Protocol: resolvedProto,
			Features: features,
		}

		fo.OutboundChannel <- completedBatch
	}
}