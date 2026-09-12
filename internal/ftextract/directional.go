package ftextract

func GetDirectionalFeatures(feats *Features, isFwd []bool, payloadLens []float32) {
	var fwdPkts, revPkts float32
	var fwdBytes, revBytes float32

	n := len(isFwd)
	for i := 0; i < n; i++ {
		plen := float32(0.0)
		if i < len(payloadLens) {
			plen = payloadLens[i]
		}

		if isFwd[i] {
			fwdPkts++
			fwdBytes += plen
		} else {
			revPkts++
			revBytes += plen
		}
	}

	feats.FwdPackets = fwdPkts
	feats.RevPackets = revPkts
	feats.FwdPayloadBytes = fwdBytes
	feats.RevPayloadBytes = revBytes

	totalPkts := fwdPkts + revPkts
	if totalPkts > 0 {
		feats.DirNormAsymPackets = (fwdPkts - revPkts) / totalPkts
	} else {
		feats.DirNormAsymPackets = 0.0
	}

	totalBytes := fwdBytes + revBytes
	if totalBytes > 0 {
		feats.DirNormAsymBytes = (fwdBytes - revBytes) / totalBytes
	} else {
		feats.DirNormAsymBytes = 0.0
	}
}