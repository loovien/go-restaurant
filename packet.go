package main

import (
	"encoding/binary"
)

// DinnerPacket
type DinnerPacket struct {
	length uint32
	body []byte
}

// Serialize
func (dinnerPacket *DinnerPacket) Serialize() []byte {
	headerLen := uint32(len(dinnerPacket.body))
	sendBody := make([]byte, headerLen + DINNER_PACKET_HEADER_LENGTH)
	binary.LittleEndian.PutUint32(sendBody[:4], headerLen)
	copy(sendBody[4:], dinnerPacket.body)
	return sendBody
}

func NewDinnerPacket(headLen uint32, body []byte) *DinnerPacket {
	return &DinnerPacket{headLen, body}
}
