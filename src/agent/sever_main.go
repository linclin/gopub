package main

import (
	"fmt"
	"github.com/linclin/gopub/src/library/p2p/common"
	"github.com/linclin/gopub/src/library/p2p/server"
	"os"
	"os/signal"
)

func main() {
	var cfg common.Config
	cfg.Name = "test"
	cfg.Auth.Password = "1234"
	cfg.Auth.Password = "gopub"
	cfg.DownDir = "/data/gcz/gopub/src/agent/testData1/"

	ss, err := common.ParserConfig(&cfg)
	cfg.Net.DataPort = 45002
	cfg.Net.MgntPort = 45003
	fmt.Print("111111111111", ss, err)
	svc, err := server.NewServer(&cfg)
	if err != nil {
		fmt.Printf("start server error, %s.\n", err.Error())
		os.Exit(4)
	}
	fmt.Print("111111111111", svc, err)
	if err = svc.Start(); err != nil {
		fmt.Printf("Start service failed, %s.\n", err.Error())
		os.Exit(4)
	}
	quitChan := listenSigInt()
	select {
	case <-quitChan:
		fmt.Printf("got control-C")
		svc.Stop()
	}
}
func listenSigInt() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	return c
}
