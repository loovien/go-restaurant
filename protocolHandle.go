package main

import (
	"github.com/vvotm/gotcp"
	log "github.com/cihub/seelog"
	"encoding/json"
	"time"
)

type RespData struct {
	Code int `json:"code"`
	Data interface{} `json:"data"`
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
//	"acceptUnion": 1 // accept eat with other at same table
// }

func (protocolHandle *ProtocolHandle) EmptySeat(conn *gotcp.Conn, param map[string]interface{}) bool {
	log.Infof("EmptySeat Handle: %v", param)

	data := map[string]interface{}{"emptyNum": "10", "tableNo": []string{"A01", "A02", "A03"}}
	byteData, err := GetRespData(0, data)
	if err != nil {
		log.Errorf("Json数据错误 %v", err)
		return false
	}
	sendPacket := NewDinnerPacket(uint32(len(byteData)), byteData)
	conn.AsyncWritePacket(sendPacket, time.Second * 3)
	return true
}

// SitDown
// packet information
// {
//	"cmd": "sitDown",
//	"tableNo": "A01", // 桌号
//	"token" : "token", // 座位口令, 正确口令才能坐下
// }
func (protocolHandle *ProtocolHandle) SitDown(conn *gotcp.Conn, param map[string]interface{}) bool {
	log.Infof("SitDown Handle: %v", param)
	tableNoInterface, okNo := param["tableNo"]
	tokenInterface, okToken := param["token"]
	if !okNo || !okToken {
		return sendPacket(conn, 1, "桌号, 身份不能为空")
	}
	tableNo := tableNoInterface.(string)
	token := tokenInterface.(string)
	tableNum := GetNumByTableNo(tableNo)

	brain, _ := GetRTBrain()
	tableInfo, ok := brain.TableConn[tableNo]
	if !ok || len(tableInfo.People) == 0 {
		brain.TableConn[tableNo] = &TableInfo{TableNo:tableNo, Token:token, People: []*gotcp.Conn{conn}}
		brain.TableNoConn[conn.GetRawConn().RemoteAddr().String()] = tableNo
		return sendPacket(conn, 0, "瓜子, 花生先吃起来!!!") // 上花生瓜子
	}

	for _, item := range tableInfo.People {
		if conn.GetRawConn().RemoteAddr().String() == item.GetRawConn().RemoteAddr().String() {
			return sendPacket(conn, 1, "您已经就坐, 还要干嘛?") // 桌子已满
		}
	}

	if len(tableInfo.People) >= tableNum {
		return sendPacket(conn, 1, "桌位已经满了, 坐不下了!") // 桌子已满
	}
	if token != tableInfo.Token {
		return sendPacket(conn, 1, "你是那来的, 不是这桌的把!")
	}
	tableInfo.People = append(tableInfo.People, conn) // join to table
	brain.TableNoConn[conn.GetRawConn().RemoteAddr().String()] = tableNo
	return sendPacket(conn, 0, "欢迎您, 朋友, 等你很久了, 瓜子花生都吃完了!")
}

// OrderMeat start order meat
// packet information
// {
//	"cmd": "orderMeat",
//	"tableNo": "A02", // 桌号
//	"token" : "token", // 座位口令, 正确口令才能坐下
//	"menu":["gbjd", "yxrs"], // 菜谱
// }
func (protocolHandle *ProtocolHandle) OrderMeat(conn *gotcp.Conn, param map[string]interface{}) bool {
	log.Infof("OrderMeat Handle:%v", param)
	tableNoInterface, okNo := param["tableNo"]
	tokenInterface, okToken := param["token"]
	if !okNo || !okToken {
		return sendPacket(conn, 1, "桌号, 身份不能为空")
	}
	tableNo := tableNoInterface.(string)
	token := tokenInterface.(string)
	brain, _ := GetRTBrain()
	tableInfo, ok := brain.TableConn[tableNo]
	if !ok || token != tableInfo.Token {
		return sendPacket(conn, 1, "你不能给别人的桌点菜!")
	}
	menu, ok := param["menu"]
	if !ok {
		return sendPacket(conn, 1, "点菜啊!")
	}
	conf, _ := GetConf()
	order := Order{TableNo:tableNo, Token:token, Food:[]Recipe{}}
	for _, recipeInterface := range menu.([]interface{})  {
		recipeName := recipeInterface.(string)
		recipe, ok := conf.Restaurant.Menu[recipeName]
		if ok {
			order.Food = append(order.Food, Recipe{}, recipe)
		}

	}
	if len(order.Food) == 0 {
		return sendPacket(conn, 1, "不要下点菜啊, 都没有哇")
	}
	brain.OrderChan <- order
	return sendPacket(conn, 0, "点菜成功, 我们马上上菜")
	return true
}



// Menu Get restaurant provide menu list
// packet information
// {
//	"cmd": "menu"
// }
func (protocolHandle *ProtocolHandle) Menu(conn *gotcp.Conn, param map[string]interface{})  bool {
	conf, isLoad := GetConf()
	if !isLoad {
		log.Info("config file not initialized!")
	}
	resp, _ := GetRespData(0, conf.Restaurant.Menu)
	menuPacket := NewDinnerPacket(uint32(len(resp)), resp)
	conn.AsyncWritePacket(menuPacket, time.Second * 3)
	return true
}

