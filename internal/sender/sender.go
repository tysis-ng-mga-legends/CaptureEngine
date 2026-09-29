package sender

import (
	"capture_engine/internal/pkcapture"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-zeromq/zmq4"
)

func StartFeatureStream(orchestrator *pkcapture.FlowOrchestrator) {
	ctx := context.Background()
	pusher := zmq4.NewPush(ctx)
	err := pusher.Dial("ipc:///tmp/flow_pipeline.ipc")
	if err != nil {
		log.Fatalf("Failed to connect to ZeroMQ pipeline: %v", err)
	}
	defer pusher.Close()

	fmt.Println("[ZMQ SENDER] Ready to transmit feature data.")

	for batch := range orchestrator.OutboundChannel {
		// Stamp precise Go departure time
		batch.Telemetry.TGoSendMicro = time.Now().UnixMicro()

		payload := map[string]any{
			"flow_id":   batch.ID,
			"protocol":  batch.Protocol,
			"features":  batch.Features,
			"telemetry": batch.Telemetry,
		}

		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Error marshalling flow data to JSON: %v", err)
			continue
		}

		msg := zmq4.NewMsgFrom(jsonBytes)
		err = pusher.Send(msg)
		if err != nil {
			log.Printf("Error pushing to ZeroMQ: %v", err)
		} else {
			fmt.Printf("[EXPORT] Flow %s:%d <-> %s:%d | Proto: %s | Twin: %d µs | Text: %d µs\n",
				batch.ID.SrcIP, batch.ID.SrcPort, batch.ID.DstIP, batch.ID.DstPort,
				batch.Protocol, batch.Telemetry.TWindowMicro, batch.Telemetry.TExtractMicro)
		}
	}
}