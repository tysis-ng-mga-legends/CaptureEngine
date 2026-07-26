package ftextract

import "fmt"

type Features struct {
	FwdPktLenMax int16
	BwdPktLenMax int16
	FwdPktLenMin int16
	BwdPktLenMin int16
	FwdPktLenMean float32
	BwdPktLenMean float32
	FwdPktLenSTD float32
	BwdPktLenSTD float32
	PktLenVar float32
	FlowIatMean float32
	FlowIatSTD float32
	FlowIatMax float32
	FwdIatMean float32
	// PshFlagCnt int8
}

func (f *Features)ExractFeatures(length []int, isFwd []bool, timestamp []int64) {
	
	tf := TemporalFeatures{}
	temp := tf.GetTemporalFeatures(timestamp)
	// fmt.Print(timestamp, length)	
	fmt.Print(temp)

}

