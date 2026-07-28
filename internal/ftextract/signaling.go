package ftextract

func ExtractSignalingFeatures(feats *Features, pshFlags []bool) {
	var pshCount uint8 = 0

	for _, psh := range pshFlags {
		if psh {
			pshCount++
		}
	}
	feats.PshFlagCount = pshCount
}