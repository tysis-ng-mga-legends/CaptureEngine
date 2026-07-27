package ftextract
import (
		"math"
)

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

func calculateMin(data []float32) int16 {
	if len(data) == 0 {
		return 0
	}
	minval := data[0]
	for _, v := range data {
		if v < minval {
			minval = v
		}
	}
	return int16(minval)
}

func calculateMax(data []float32) int16 {
	if len(data) == 0 {
		return 0
	}
	maxval := data[0]
	for _, v := range data {
		if v > maxval {
			maxval = v
		}
	}
	return int16(maxval)
}

func calculateMean(data []float32) float32 {
	if len(data) == 0 {
		return 0
	}
	var sum float32
	for _, v := range data {
		sum += v
	}
	return float32(sum / float32(len(data)))
}

func calculateVariance(data []float32) float32 {
	n := len(data)
	if n <= 1 {
		return 0
	}
	mu := calculateMean(data)
	var sumSquares float32
	for _, v := range data {
		sumSquares += (v - mu) * (v - mu)
	}
	return float32(sumSquares / float32(n - 1))
}

func calculateStdDev(data []float32) float32 {
	return float32(math.Sqrt(float64(calculateVariance(data))))
}
