package ftextract

import (
	"math"
)

type Features struct {

	FrameLenMean float32
	FrameLenStd float32
	FrameLenMin float32
	FrameLenMax float32
	FrameLenCV float32

	PayloadLenMean float32
	PayloadLenStd float32
	PayloadLenMin float32
	PayloadLenMax float32
	PayloadLenSum float32
	PayloadLenCV float32

	Duration float32
	IatMean float32
	IatStd float32
	IatMin float32
	IatMax float32
	IatCV float32

	FwdPackets float32
	RevPackets float32
	FwdPayloadBytes float32
	RevPayloadBytes float32
	DirNormAsymPackets float32
	DirNormAsymBytes float32
}

func ExtractFeatures(frameLens []float32, payloadLens []float32, isFwd []bool, timestamp []int64, pshFlags []bool, proto string) Features {
	if len(frameLens) == 0 {
		return Features{}
	}
	var feats Features

	frameFeats := ExtractFrameFeatures(frameLens)
	feats.FrameLenMean = frameFeats.Mean
	feats.FrameLenStd = frameFeats.Std
	feats.FrameLenMin = frameFeats.Min
	feats.FrameLenMax = frameFeats.Max
	feats.FrameLenCV = frameFeats.CV

	payloadStats := CalcStats(payloadLens)
	feats.PayloadLenMean = payloadStats.Mean
	feats.PayloadLenStd = payloadStats.Std
	feats.PayloadLenMin = payloadStats.Min
	feats.PayloadLenMax = payloadStats.Max
	feats.PayloadLenSum = payloadStats.Sum
	feats.PayloadLenCV = payloadStats.CV

	GetTemporalFeatures(&feats, timestamp)
	GetDirectionalFeatures(&feats, isFwd, payloadLens)
	ExtractSignalingFeatures(&feats, pshFlags, proto)

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


