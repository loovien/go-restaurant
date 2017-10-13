package main

import (
	"syscall"
	"os/signal"
	"fmt"
	"os"
	log "github.com/cihub/seelog"
)

func main() {
	Bootstrap()

	log.Debug("ooooooooooooooooooooooooooooooooo")

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	fmt.Println("System Signal:", <-signalChan)
}

func Bootstrap()  {
	InitLog("conf/log4g.xml")
	InitConf("conf/app.toml")
}
