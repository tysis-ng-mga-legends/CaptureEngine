package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"capture_engine/internal/pkcapture"

	"github.com/go-zeromq/zmq4"
)

func StartFeatureStream(orchestrator *pkcapture.FlowOrchestrator) {

	ctx := context.Background()

	pusher := zmq4.NewPush(ctx)
	
	err := pusher.Dial("ipc:///tmp/flow_pipeline.ipc")

	if err != nil{
		log.Fatalf("Failed to connect to ZeroMQ pipeline: %v", err)
	}
	
	defer pusher.Close()
	fmt.Println("Sending Feature Data")

	for batch := range orchestrator.OutboundChannel {
		
		payload := map[string]any{
			"flow_id": batch.ID,
			"protocol": batch.Protocol,
			"features": batch.Features,
		}


		jsonBytes, err := json.Marshal(payload)
		if err != nil{
			log.Printf("Error marshalling flow data to JSON: %v", err)
			continue
		}

		msg := zmq4.NewMsgFrom(jsonBytes)
		err = pusher.Send(msg)

		if err != nil{
			log.Printf("Error pushing to ZeroMQ: %v", err)
		} else {
			fmt.Printf("[ZMQ SENDER] Exported flow window for %s:%d -> %s\n", batch.ID.SrcIP, batch.ID.SrcPort, batch.Protocol)
		}
		
	}  
		
}