package pkcapture


type FlowID struct {
	SrcIP  string
	DstIP string
	SrcPort uint16
	DstPort uint16
	Protocol string
}


func (id FlowID) GetNormalized() FlowID {
	if id.SrcIP > id.DstIP {
		return FlowID{
			SrcIP:    id.DstIP,
			DstIP:    id.SrcIP,
			SrcPort:  id.DstPort,
			DstPort:  id.SrcPort,
			Protocol: id.Protocol,
		}
	}

	if id.SrcIP == id.DstIP && id.SrcPort > id.DstPort {
		return FlowID{
			SrcIP:    id.SrcIP,
			DstIP:    id.DstIP,
			SrcPort:  id.DstPort,
			DstPort:  id.SrcPort,
			Protocol: id.Protocol,
		}
	}

	return id
}