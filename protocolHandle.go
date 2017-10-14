package main

import (
	"github.com/vvotm/gotcp"
	log "github.com/cihub/seelog"
	"encoding/json"
	"time"
)

type RespData struct {
	Code int
	Data interface{}
}

func GetRespData(code int, data interface{}) ([]byte, error){
	respData := RespData{Code:code, Data:data}
	return json.Marshal(respData)
}

type ProtocolHandle struct {

}

// have empty seat or not
// ask packet information
// {
// 	"cmd": "emptyseat",
//	"leader": "boos", // who pay order
// 	"eatNum": 12, // number of eat
//	"acceptUnion": 1 // accept eat with other at same table
// }

func (protocolHandle *ProtocolHandle) EmptySeat(conn *gotcp.Conn, param map[string]interface{}) bool {
	log.Infof("EmptySeat Handle: %v", param)

	data := map[string]interface{}{"emptyNum": "10"}
	byteData, err := GetRespData(0, data)
	if err != nil {
		log.Errorf("Json数据错误 %v", err)
		return false
	}
	sendPacket := NewDinnerPacket(uint32(len(byteData)), byteData)
	conn.AsyncWritePacket(sendPacket, time.Second * 3)
	return true
}
