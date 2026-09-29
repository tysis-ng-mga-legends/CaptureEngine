package pkcapture

import (
	"capture_engine/internal/ftextract"
	"sync"
	"time"
)

type PacketData struct {
	SourceIP   string `json:"source_ip"`
	SourcePort uint16 `json:"source_port"`
	Protocol   string `json:"proto"`
	Timestamp  int64  `json:"timestamp"` // Microseconds
	Length     int    `json:"length"`
	PayloadLen int    `json:"payload_len"`
}

type TelemetryData struct {
	TWindowMicro  int64 `json:"t_window_us"`
	TExtractMicro int64 `json:"t_extract_us"`
	TGoSendMicro  int64 `json:"t_go_send_us"`
}

type FlowBatch struct {
	ID        FlowID             `json:"flow_id"`
	Protocol  string             `json:"protocol"`
	Features  ftextract.Features `json:"features"`
	Telemetry TelemetryData      `json:"telemetry"`
}

type FlowOrchestrator struct {
	mu             sync.Mutex
	ActiveSequence map[FlowID][]PacketData
	ProcessedFlows map[FlowID]bool
	OutboundChannel chan FlowBatch
}

func NewOrchestrator() *FlowOrchestrator {
	return &FlowOrchestrator{
		ActiveSequence:  make(map[FlowID][]PacketData),
		ProcessedFlows:  make(map[FlowID]bool),
		OutboundChannel: make(chan FlowBatch, 100),
	}
}

func (fo *FlowOrchestrator) IncrementPacketCount(flowID FlowID, packet PacketData) {
	canonicalID := flowID.GetNormalized()

	fo.mu.Lock()
	defer fo.mu.Unlock()

	// Once a flow window is completed, skip further packets for it
	if fo.ProcessedFlows[canonicalID] {
		return
	}

	fo.ActiveSequence[canonicalID] = append(fo.ActiveSequence[canonicalID], packet)

	// Exactly 7 packets trigger extraction
	if len(fo.ActiveSequence[canonicalID]) == 8 {
		packetBatch := fo.ActiveSequence[canonicalID]
		fo.ProcessedFlows[canonicalID] = true
		delete(fo.ActiveSequence, canonicalID) // Free memory immediately

		startExtract := time.Now()

		resolvedProto := canonicalID.Protocol
		for _, p := range packetBatch {
			if p.Protocol == "TLS" || p.Protocol == "QUIC" {
				resolvedProto = p.Protocol
				break
			}
		}

		frameLens := make([]float32, 8)
		payloadLens := make([]float32, 8)
		isFwd := make([]bool, 8)
		timestamps := make([]int64, 8)

		initiatorIP := packetBatch[0].SourceIP
		initiatorPort := packetBatch[0].SourcePort

		for i, pkt := range packetBatch {
			frameLens[i] = float32(pkt.Length)
			payloadLens[i] = float32(pkt.PayloadLen)
			isFwd[i] = (pkt.SourceIP == initiatorIP && pkt.SourcePort == initiatorPort)
			timestamps[i] = pkt.Timestamp
		}

		features := ftextract.ExtractFeatures(frameLens, payloadLens, isFwd, timestamps)
		extractDur := time.Since(startExtract).Microseconds()
		windowDur := timestamps[7] - timestamps[0]

		completedBatch := FlowBatch{
			ID:       canonicalID,
			Protocol: resolvedProto,
			Features: features,
			Telemetry: TelemetryData{
				TWindowMicro:  windowDur,
				TExtractMicro: extractDur,
			},
		}

		fo.OutboundChannel <- completedBatch
	}
}