package init_sever

import (
	"fmt"
	"github.com/astaxie/beego"
	"library/p2p/common"
	"library/p2p/server"
	"os"
)

var P2pSvc *server.Server

func init() {
	cfg := common.ReadJson("agent/server.json")
	ss, err := common.ParserConfig(&cfg)
	cfg.Server = true
	fmt.Print("ctg:", ss, err)
	P2pSvc, err = server.NewServer(&cfg)
	if err != nil {
		beego.Error("start server error, %s.\n", err.Error())
		if beego.BConfig.RunMode != "docker" {
			os.Exit(4)
		}
	}
}
func Start() {
	if err := P2pSvc.Start(); err != nil {
		beego.Error("Start service failed, %s.\n", err.Error())
		if beego.BConfig.RunMode != "docker" {
			os.Exit(4)
		}
	}
}
