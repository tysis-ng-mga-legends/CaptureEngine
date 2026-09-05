package ftextract

import (
	"math"
)

type FrameFeatures struct {
	Mean float64
	Std  float64
	Min  float64
	Max  float64
	CV   float64
}

// ExtractFrameFeatures computes features 1-5 across all 7 packets
func ExtractFrameFeatures(frameLens []float64) FrameFeatures {
	n := len(frameLens)
	if n == 0 {
		return FrameFeatures{}
	}

	minVal := frameLens[0]
	maxVal := frameLens[0]
	for _, v := range frameLens {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
		sumVal += v
	}

	meanVal := sumVal / float64(n)

	// ddof = 1 match np.std(ddof-1) for sample stddev
	stdVal := 0.0
	if n >= 2 {
		var sumSquares float64
		for _, v := range frameLens {
			diff := v - meanVal
			sumSquares += diff * diff
		}
		stdVal = math.Sqrt(sumSquares / float64(n-1))
	}

	cvVal := 0.0
	if meanVal != 0 {
		cvVal = stdVal / meanVal
	}

	return FrameFeatures{
		Mean: meanVal,
		Std:  stdVal,
		Min:  minVal,
		Max:  maxVal,
		CV:   cvVal,
	}
}
