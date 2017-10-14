package main

import (
	"testing"
	log "github.com/cihub/seelog"
)

func TestLog(t *testing.T)  {
	InitLog("")
	log.Debug("-------------")
}
