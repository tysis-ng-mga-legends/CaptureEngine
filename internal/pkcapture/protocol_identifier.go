package pkcapture

func isTLS(payload []byte) bool {
	if (len(payload) < 5) {
		return false
	}

	contentType := payload[0]

	if contentType < 0x14  || contentType > 0x17 {
		return false
	}

	version := uint16(payload[1]) << 8 | uint16(payload[2])

	if version < 0x0301 || version > 0x0304 {
		return false
	}

	recordLen := int(payload[3]) << 8 | int(payload[4])

	return recordLen > 0 && recordLen < 18432
}

