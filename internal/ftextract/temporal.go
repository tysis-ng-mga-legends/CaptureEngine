package ftextract

import (
	"fmt"
)


type TemporalFeatures struct {
	FlowIatMean float32
	FlowIatSTD float32
	FlowIatMax float32
	FwdIatMean float32
}

func (tf *TemporalFeatures) GetTemporalFeatures(timestamp []uint64){
	
	// for _ ,data := range packetData{
	// 	fmt.Printf(`Timestamp: %s, Length: %d`+"\n", data.Timestamp.Format(time.RFC3339), data.Length)
	// }

	fmt.Print(timestamp)
}

// func getFlowIAtMean ()