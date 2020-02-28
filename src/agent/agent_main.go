package main

import (
	"fmt"
	"github.com/linclin/gopub/src/library/p2p/agent"
	"github.com/linclin/gopub/src/library/p2p/common"
	"os"
	"os/signal"
)

func main() {
	cfg := common.ReadJson("agent.json")
	ss, err := common.ParserConfig(&cfg)
	fmt.Print("Config:", ss)
	svc, err := agent.NewAgent(&cfg)
	if err != nil {
		fmt.Printf("start agent error, %s.\n", err.Error())
		os.Exit(4)
	}
	if err = svc.Start(); err != nil {
		fmt.Printf("Start service failed, %s.\n", err.Error())
		os.Exit(4)
	}
	quitChan := listenSigInt1()
	select {
	case <-quitChan:
		fmt.Printf("got control-C")
		svc.Stop()
	}
}
func listenSigInt1() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	return c
}
