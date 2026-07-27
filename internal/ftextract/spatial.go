package ftextract

type SpatialFeatures struct {
	FwdPktLenMax int16
	BwdPktLenMax int16
	FwdPktLenMin int16
	BwdPktLenMin int16
	FwdPktLenMean float32
	BwdPktLenMean float32
	FwdPktLenSTD float32
	BwdPktLenSTD float32
	PktLenVar float32
}

// ExtractSpatialFeatures parses the window and determines packet direction 
func ExtractSpatialFeatures(lengths []int, isFwd []bool) SpatialFeatures {
	var feats SpatialFeatures
	
	if len(lengths) == 0 {
		return feats
	}

	var lFwd []float32
	var lBwd []float32
	var lAll []float32

	for i, length := range lengths {
		length := float32(length)
		
		// all lengths are stored for bidirectional calculations
		lAll = append(lAll, length)

		// packets are routed based on their direction relative to the initiator
		if isFwd[i] {
			lFwd = append(lFwd, length)
		} else {
			lBwd = append(lBwd, length)
		}
	}

	feats.FwdPktLenMin = calculateMin(lFwd)
	feats.BwdPktLenMin = calculateMin(lBwd)
	feats.FwdPktLenMax = calculateMax(lFwd)
	feats.BwdPktLenMax = calculateMax(lBwd)
	feats.FwdPktLenMean = calculateMean(lFwd)
	feats.BwdPktLenMean = calculateMean(lBwd)
	feats.FwdPktLenSTD = calculateStdDev(lFwd)
	feats.BwdPktLenSTD = calculateStdDev(lBwd)
	feats.PktLenVar = calculateVariance(lAll)
	
	return feats
}
