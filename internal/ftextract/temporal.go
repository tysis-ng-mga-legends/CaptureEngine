package ftextract

import (
	"fmt"
	"time"
)

type PacketInput struct {
	Timestamp time.Time
	Length int
}
type TemporalFeatures struct {
	FlowIatMean float32
	FlowIatSTD float32
	FlowIatMax float32
	FwdIatMean float32
}

func (tf *TemporalFeatures) GetTemporalFeatures(packetData []PacketInput){
	
	// for _ ,data := range packetData{
	// 	fmt.Printf(`Timestamp: %s, Length: %d`+"\n", data.Timestamp.Format(time.RFC3339), data.Length)
	// }

	fmt.Print(packetData)
}

// func getFlowIAtMean ()