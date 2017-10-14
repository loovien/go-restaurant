package main

import (
	"syscall"
	"os/signal"
	"fmt"
	"os"
	log "github.com/cihub/seelog"
	"github.com/vvotm/gotcp"
	"flag"
	"net"
	"time"
)

var listenAddr *string = flag.String("addr", "", "address of listen")
func main() {
	Bootstrap()
	Run()
}

func Bootstrap()  {
	InitLog("conf/log4g.xml")
	InitConf("conf/app.toml")
	conf, _ := GetConf()
	InitTRBran(conf.Restaurant.CookNum)
}

func Run()  {
	log.Info("*********Server Start**********")
	tcpconf := &gotcp.Config{
		PacketSendChanLimit: 10, // the limit of packet send channel
		PacketReceiveChanLimit: 10, // the limit of packet receive channel
	}
	if *listenAddr == "" { // if command line not pass address use config file
		conf, _ := GetConf()
		listenAddr = &conf.Addr
	}
	log.Infof("Listen at: %s", *listenAddr)
	srv := gotcp.NewServer(tcpconf, &Callback{}, &DinnerProtocol{})
	tcpAddr, _ := net.ResolveTCPAddr("tcp", *listenAddr)
	tcpListenAddr, err := net.ListenTCP("tcp", tcpAddr)
	defer tcpListenAddr.Close()
	if err != nil {
		log.Errorf("Listen: %v", err)
		os.Exit(1)
	}
	go srv.Start(tcpListenAddr, time.Second * 10)

	brain, err := GetRTBrain()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	brain.CookWork() // 厨师上班
	brain.RecipeServing() // 上菜上班
	log.Info("********Restaurant Cook Working ********")
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	log.Info("********ALL OK********")
	fmt.Println("System Signal:", <-signalChan)
	brain.Stop()
	srv.Stop()
}
