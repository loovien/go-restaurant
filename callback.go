package main

import (
	"github.com/vvotm/gotcp"
	log "github.com/cihub/seelog"
	"encoding/json"
	"reflect"
	"strings"
)

type Callback struct {

}

// OnConnect when tcp connect this method will be Invoke
func (callback *Callback) OnConnect(conn *gotcp.Conn) bool {
	addr := conn.GetRawConn().RemoteAddr().String()
	log.Infof("%s connected", addr)
	conn.PutExtraData(addr)
	return true
}

// OnMessage when establish connection data communicate this method will be invoke
func (callback *Callback) OnMessage(conn *gotcp.Conn, packet gotcp.Packet) bool {
	receivePacket  := packet.(*DinnerPacket)
	if len(receivePacket.body) == 0 {
		return false
	}
	protocolData := make(map[string]interface{})
	log.Infof("收到客户端数据 : %v", string(receivePacket.body))
	err := json.Unmarshal(receivePacket.body, &protocolData)
	if err != nil {
		log.Errorf("反JSON数据错误: %v", err)
		return false
	}
	cmdid, ok := protocolData["cmd"]
	if !ok {
		log.Errorf("不能识别的协议")
		return false
	}
	protocolHandles := &ProtocolHandle{}
	protocolRefValue := reflect.ValueOf(protocolHandles)
	methodRefValue := protocolRefValue.MethodByName(strings.Title(cmdid.(string)))
	if !methodRefValue.IsValid() {
		log.Errorf("协议对应的处理不存在")
		return false
	}
	methodRefValue.Call([]reflect.Value{reflect.ValueOf(conn), reflect.ValueOf(protocolData)})
	return true
}

// OnClose when establish connection close this method will be invoke
func (callback *Callback) OnClose(conn *gotcp.Conn) {
	log.Infof("%v Closed!", conn.GetExtraData())
}
