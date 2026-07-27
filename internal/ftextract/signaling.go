package ftextract

func ExtractSignalingFeatures(feats *Features, tcpFlags []uint8) {
	var pshCount float32 = 0

	for _, flag := range tcpFlags {
		if (flagByte & 0x08) != 0 {
			pshCount++
		}
	}
	feats.PshFlagCount = pshCount
}