package ftextract

import (
	"sort"
)

// function that coordinates the functions in this file and returns TemporalFeatures
func GetTemporalFeatures(feats *Features, timestamps []int64) {
	if len(timestamps) < 2 {
		feats.Duration = 0
		feats.IatMean = 0
		feats.IatStd = 0
		feats.IatMin = 0
		feats.IatMax = 0
		feats.IatCV = 0
		return
	}

	// Sort timestamps to ensure they are in ascending order
	sorted := make([]int64, len(timestamps))
	copy(sorted, timestamps)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	// Duration (t_last - t_first) in seconds
	feats.Duration = float32(sorted[len(sorted)-1]-sorted[0]) / 1e6

	// Compute the 6 iat in seconds
	iats := make([]float32, len(sorted)-1)
	for i := 1; i < len(sorted); i++ {
		delta := sorted[i] - sorted[i-1]
		if delta < 0 {
			delta = 0
		}
		iats[i-1] = float32(delta) / 1e6
	}

	// 13 - 17 Statistics of IATs
	stats := CalcStats(iats)
	feats.IatMean = stats.Mean
	feats.IatStd = stats.Std
	feats.IatMin = stats.Min
	feats.IatMax = stats.Max
	feats.IatCV = stats.CV
}