package ftextract

import (
	"sort"
)

// function that coordinates the functions in this file and returns TemporalFeatures
func GetTemporalFeatures(feats *Features, timestamps []int64, isFwd []bool) {

	iats := computeIat(timestamps)

	feats.FlowIatMean = calculateMean(iats)
	feats.FlowIatSTD = calculateStdDev(iats)
	feats.FlowIatMax = findIatMax(iats)

	var fwdTimestamps []int64
	for i, timestamp := range timestamps {
		if isFwd[i]{
			fwdTimestamps = append(fwdTimestamps, timestamp)
		}
	}

	fwdIats := computeIat(fwdTimestamps)
	feats.FwdIatMean = calculateMean(fwdIats)
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

// helper function that finds the max
func findIatMax(values []float32) float32 {
	max := values[0]
	for _, v := range values{
		if v > max{
			max = v
		}
	}
	return max
}
