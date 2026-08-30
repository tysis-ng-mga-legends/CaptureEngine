package pkcapture

import "github.com/google/gopacket/layers"

// 1. Structural Invariant Dissectors

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

	return recordLen > 0 && recordLen <= 18432
}

func isQUIC(payload []byte) bool {

	if len(payload) < 1 {
		return false
	}

	firstByte := payload[0]

	if (firstByte & 0x40) == 0 {
		return false
	}

	isLongHeader := (firstByte & 0x80) != 0

	if isLongHeader {
		if len(payload) < 5 {
			return false
		}

		version :=  uint32(payload[1]) << 24 | uint32(payload[2]) << 16 | uint32(payload[3]) << 8 | uint32(payload[4])

		return version == 0x00000001 || version == 0x6b3343cf || (version & 0xff000000) == 0xff000000 ||
		       (version & 0xff000000) == 0x51000000

	}

	return len(payload) >= 21
}

// 2. Protocol-Specific Application Data Gatekeepers
func isTLSAppData(payload []byte) bool {
	if !isTLS(payload) {
		return false
	}

	return payload[0] == 0x17
}

func isPlainTCPAppData(tcp *layers.TCP) bool {
	if tcp.SYN || tcp.ACK || tcp.FIN || len(tcp.Payload) == 0 {
		return false
	}

	return true
}

func isQUICAppData(payload []byte) bool {
	if !isQUIC(payload){
		return false
	}

	return (payload[0] & 0x80) == 0
}

