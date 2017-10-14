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
				default:
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
			case order := <- rtBrain.ServingChan:
				tableInfo, ok := rtBrain.TableConn[order.TableNo]
				if !ok {
					return
				}
				for _, people := range tableInfo.People { // 上菜
					recipeByte, _ := json.Marshal(order.Food)
					sendPacket(people, 0, string(recipeByte))
					return
				}
			}
		}
	}()
}

// Stop 关门
func (rtBrain *RTBrain) Stop() {
	close(rtBrain.OrderChan)
	close(rtBrain.ServingChan)
	close(rtBrain.CloseChan)
	log.Info("Wait ALL Order Completed!")
	rtBrain.waitGroup.Wait()
}
