package main

import (
	log "github.com/cihub/seelog"
	"github.com/BurntSushi/toml"
)

var confInstance *Conf

type Recipe struct {
	Name string `json:"name"`
	Price float64 `json:"price"`
	Stock int `json:"stock"`
}

type Restaurant struct {
	OpenAt string `json:"openAt"`
	CloseAt string `json:"closeAt"`
	TotalSeat int `json:"totalSeat"`
	Table2p int `json:"table2p"`
	Table4p int `json:"table4p"`
	Table6p int `json:"table6p"`
	Table10p int `json:"table10p"`
	CookNum int `json:"cookNum"`
	Menu map[string] Recipe `json:"menu"`
}

type Conf struct{
	Name string
	Version string
	Addr string
	Restaurant Restaurant
}

// NewConf
func NewConf() *Conf {
	return &Conf{}
}

// InitConf
func InitConf(filename string)  {
	NewConf().Load(filename, true)
}

// GetConf get conf instance
func GetConf() (*Conf, bool) {
	if confInstance == nil {
		return NewConf(), false
	}
	return confInstance, true
}

// Load loading config
func (conf *Conf) Load (filename string, reload bool) {
	if confInstance != nil && !reload {
		return
	}
	configFile, err := getAbsPath(filename)
	if err != nil {
		log.Infof("file 【%s】 not exists use default file conf/app.toml", filename)
		configFile = DEFAULT_CONF_FILE
	}

	tmpConf := NewConf()
	_, err = toml.DecodeFile(configFile, tmpConf)
	if err != nil {
		log.Criticalf("config file 【%s】 %v !", configFile, err)
	}
	confInstance = tmpConf
}