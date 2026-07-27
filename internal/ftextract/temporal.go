package ftextract

import (
	"math"
	"sort"
)


type TemporalFeatures struct {
	FlowIatMean float32
	FlowIatSTD float32
	FlowIatMax float32
	FwdIatMean float32
}

// function that coordinates the functions in this file and returns TemporalFeatures
func (tf *TemporalFeatures) GetTemporalFeatures(timestamps []int64) TemporalFeatures{

	iats := computeIat(timestamps)
	iatMean := calculateMean(iats)
	iatStd := calculateStd(iatMean, iats)
	iatMax := findMax(iats)

	return TemporalFeatures{
		FlowIatMean: iatMean,
		FlowIatSTD: iatStd,
		FlowIatMax: iatMax,
		FwdIatMean: 0.0, // still waiting for dan to finish with the directions
	}
}

// compute the inter-arrival-time of each timestamps and return an array of iats
func computeIat(timestamps []int64) []float32 {
	sorted := make([]int64, len(timestamps))
	copy(sorted, timestamps)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	iats := make([]float32, len(sorted)-1) 
	for i := 1; i < len(sorted); i++ {
		delta := sorted[i] - sorted[i-1]
		// converts microsecond to second
		iats[i-1] = float32(delta) / 1e6
	}

	return iats
}


// helper function that calculates the std
func calculateStd(mean float32, values []float32) float32{
	var varianceSum float32 = 0

	for _, t := range values{
		diff := float32(t) - mean
		varianceSum += diff * diff
	}

	return float32(math.Sqrt(float64(varianceSum/float32(len(values)))))
}

// helper function that finds the max
func findMax(values []float32) float32 {
	
	max := values[0]

	for _, v := range values{
		if v > max{
			max = v
		}
	}

	return max
}