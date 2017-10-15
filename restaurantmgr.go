package main

import (
	"sync"
	"errors"
	log "github.com/cihub/seelog"
	"github.com/vvotm/gotcp"
	"encoding/json"
)

var RTBrainInstance *RTBrain
var RTBrainOnce *sync.Once = new(sync.Once)

type Order struct {
	TableNo string
	Token string
	Food []Recipe
}

type RTBrain struct {
	RestaurantLock *sync.RWMutex
	CookNum int
	OrderChan chan Order
	ServingChan chan Order
	waitGroup *sync.WaitGroup
	TableConn map[string]*TableInfo
	CloseChan chan bool
	TableNoConn map[string]string
}

type TableInfo struct {
	TableNo string
	Token string
	People [] *gotcp.Conn
}

func InitTRBran(numCook int) *RTBrain {
	if RTBrainInstance == nil {
		RTBrainOnce.Do(func() {
			RTBrainInstance = &RTBrain{
				&sync.RWMutex{},
				numCook,
				make(chan Order, 1),
				make(chan Order, 1),
				&sync.WaitGroup{},
				make(map[string]*TableInfo),
				make(chan bool),
				make(map[string]string),
			}
		})
	}
	return RTBrainInstance
}

func GetRTBrain() (*RTBrain, error){
	if RTBrainInstance == nil {
		return nil, errors.New("RTBrain Not Initialize")
	}
	return RTBrainInstance, nil
}

// CookWork 厨师工作
func (rtBrain *RTBrain) CookWork() {
	for i := 0 ; i < rtBrain.CookNum ; i++ {
		rtBrain.waitGroup.Add(1)
		go func() {
			defer rtBrain.waitGroup.Done()
			for  {
				select {
				case <- rtBrain.CloseChan:
					return
				case order, isNotClose := <- rtBrain.OrderChan:
					if !isNotClose {
						return
					}
					log.Infof("cook food: %v", order)
					rtBrain.ServingChan <- order
				default:
					break
				}
			}
		}()
	}
}

// RecipeServing 上菜
func (rtBrain *RTBrain) RecipeServing() {
	rtBrain.waitGroup.Add(1)
	go func() {
		defer rtBrain.waitGroup.Done()
		for {
			select {
			case <- rtBrain.CloseChan:
				return
			case order, isNotClose := <- rtBrain.ServingChan:
				if !isNotClose {
					return
				}
				log.Infof("上菜: %v", order)
				tableInfo, ok := rtBrain.TableConn[order.TableNo]
				if !ok {
					log.Errorf("没吃就走了? %v", tableInfo)
					break
				}
				log.Infof("table Info: %v", tableInfo)
				log.Infof("table people Info: %v", tableInfo.People)
				log.Infof("table tableNo Info: %v", tableInfo.TableNo)
				for _, people := range tableInfo.People { // 上菜
					recipeByte, _ := json.Marshal(order.Food)
					log.Infof("send data %s", string(recipeByte))
					log.Info(people.GetExtraData())
					sendPacket(people, 0, string(recipeByte))
				}
				break
			default:
				break
			}
		}
	}()
}

// Stop 关门
func (rtBrain *RTBrain) Stop() {
	close(rtBrain.OrderChan)
	close(rtBrain.ServingChan)
	close(rtBrain.CloseChan)
	log.Info("*****Gracefully Stop*****")
	rtBrain.waitGroup.Wait()
}
