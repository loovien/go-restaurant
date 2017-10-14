package main

import (
	"github.com/vvotm/gotcp"
	"net"
	"encoding/binary"
	log "github.com/cihub/seelog"
	"errors"
)

type DinnerProtocol struct {
}
func (dinnerProtocol *DinnerProtocol) ReadPacket(conn *net.TCPConn) (gotcp.Packet, error) {
	headerByte := make([]byte, DINNER_PACKET_HEADER_LENGTH)
	conn.Read(headerByte)

	headerLen := binary.LittleEndian.Uint32(headerByte)
	if headerLen > 1024 {
		log.Criticalf("包体过长, 超过了1024字节")
		return nil, errors.New("包体过长, 超过了1024字节")
	}
	bodyByte := make([]byte, headerLen)
	conn.Read(bodyByte)
	return NewDinnerPacket(headerLen, bodyByte), nil
}
