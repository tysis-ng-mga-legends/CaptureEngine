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

func (tf *TemporalFeatures) GetTemporalFeatures(timestamps []int64){
	
	// for _ ,data := range packetData{
	// 	fmt.Printf(`Timestamp: %s, Length: %d`+"\n", data.Timestamp.Format(time.RFC3339), data.Length)
	// }

	fmt.Print(timestamps)
}

func getFlowIAtMean(timestamps []int64) float32 {
	
	var sum float32 = 0
	for _, t := range timestamps {
		sum += float32(t) 
	}
	return sum / float32(len(timestamps))
}

func getFlowIatSTD(timestamps []int64)