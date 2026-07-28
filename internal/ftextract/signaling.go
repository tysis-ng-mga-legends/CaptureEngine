package ftextract


func ExtractSignalingFeatures(feats *Features, pshFlags []bool, proto string) {
	var pshCount uint8 = 0

	for _, psh := range pshFlags {
		if psh {
			pshCount++
		}
	}
	// fmt.Println(proto)
	feats.PshFlagCount = pshCount
	feats.ProtoUDP = isUDP(proto)
}

func isUDP(proto string) uint8 {

	if proto == "UDP" || proto == "QUIC" {
		return 1
	}
	return 0
}
