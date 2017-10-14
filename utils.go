package main

import (
	"strings"
	"time"
	"github.com/vvotm/gotcp"
)

func GetNumByTableNo(tableNo string) int {
	tableNum := 0;
	prefix := tableNo[:1]
	switch strings.ToUpper(prefix) {
	case "A":
		tableNum = 2
	case "B":
		tableNum = 4
	case "C":
		tableNum = 6
	case "D":
		tableNum = 10
	}
	return tableNum
}

// sendPacket
func sendPacket(conn *gotcp.Conn, code int, msg string) bool {
	respData := map[string]string{"msg": msg}
	respByteData, _ := GetRespData(code, respData)
	missParamPacket := NewDinnerPacket(uint32(len(respByteData)), respByteData)
	conn.AsyncWritePacket(missParamPacket, time.Second * 3)
	return true
}
