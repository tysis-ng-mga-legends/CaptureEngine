package ftextract

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

func (f *Features)ExractFeatures(lengths []int, isFwd []bool, timestamps []uint64) {
	// temporal := TemporalFeatures{}
	// temporal.GetTemporalFeatures(timestamps)
	if len(lengths) == 0 {
		return
	}

	spatial := ExtractSpatialFeatures(lengths, isFwd)

	f.FwdPktLenMax = spatial.FwdPktLenMax
	f.BwdPktLenMax = spatial.BwdPktLenMax
	f.FwdPktLenMin = spatial.FwdPktLenMin
	f.BwdPktLenMin = spatial.BwdPktLenMin
	f.FwdPktLenMean = spatial.FwdPktLenMean
	f.BwdPktLenMean = spatial.BwdPktLenMean
	f.FwdPktLenSTD = spatial.FwdPktLenSTD
	f.BwdPktLenSTD = spatial.BwdPktLenSTD
	f.PktLenVar = spatial.PktLenVar
}

