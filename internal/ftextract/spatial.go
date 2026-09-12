package ftextract

import (
	"math"
)

// StatsResult holds descriptive statistics for 1D numerical feature sets
type StatsResult struct {
	Mean float32
	Std  float32
	Min  float32
	Max  float32
	Sum  float32
	CV   float32
}

// CalcStats computes mean, sample std (ddof=1), min, max, sum, and cv across a float32 slice
func CalcStats(data []float32) StatsResult {
	n := len(data)
	if n == 0 {
		return StatsResult{}
	}

	minVal := data[0]
	maxVal := data[0]
	sumVal := float32(0.0)

	for _, v := range data {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
		sumVal += v
	}

	meanVal := sumVal / float32(n)

	// Sample standard deviation (ddof=1 to match Python numpy.std(ddof=1))
	stdVal := float32(0.0)
	if n >= 2 {
		var sumSquares float32
		for _, v := range data {
			diff := v - meanVal
			sumSquares += diff * diff
		}
		stdVal = float32(math.Sqrt(float64(sumSquares / float32(n-1))))
	}

	cvVal := float32(0.0)
	if meanVal != 0 {
		cvVal = stdVal / meanVal
	}

	return StatsResult{
		Mean: meanVal,
		Std:  stdVal,
		Min:  minVal,
		Max:  maxVal,
		Sum:  sumVal,
		CV:   cvVal,
	}
}

// ExtractFrameFeatures provides backwards compatibility for existing calls
func ExtractFrameFeatures(frameLens []float32) StatsResult {
	return CalcStats(frameLens)
}