package models

type PacketData struct {
    SourceIP  string `json:"source_ip"`
    Timestamp uint64 `json:"timestamp"`
    Length    int    `json:"length"`
}