package main

import (
	"sync"
	log "github.com/cihub/seelog"
	"github.com/BurntSushi/toml"
)

var (
	confLoadOnce *sync.Once = new(sync.Once)
	confInstance *Conf
)

type Conf struct{
	Name string
	Version string
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
		log.Criticalf("Load config file 【%s】 error !", configFile)
	}
	confInstance = tmpConf
}