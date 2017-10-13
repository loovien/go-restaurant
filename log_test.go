package main

import (
	log "github.com/cihub/seelog"
	"testing"
)

func TestLog(t *testing.T)  {
	InitLog("conf/log4g.xml")
	log.Info("hello")
	log.Debug("wold")
	t.Log("------------")
}
