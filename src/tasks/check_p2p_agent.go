package tasks

import (
	"encoding/json"
	"fmt"
	"library/common"
	"library/components"
	"library/p2p/init_sever"
	"models"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils"
)

type emailConfig struct {
	UserName string `json:"username,omitempty"`
	PassWord string `json:"password,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
}

func init() {

}
func Check_p2p_angent_status() error {
	beego.Info("p2p agent自动任务开始" + time.Now().Format("2006-01-02 15:04:05"))
	var projects []models.Project
	o := orm.NewOrm()
	num, err := o.Raw("SELECT * FROM `project` WHERE  p2p=1").QueryRows(&projects)
	if num > 0 && err != nil {
		for _, project := range projects {
			s := components.BaseComponents{}
			s.SetProject(&project)
			s.SetTask(&models.Task{Id: -10})
			ips := s.GetHosts()
			ss := init_sever.P2pSvc.CheckAllClient(ips)
			reIps := []string{}
			for ip, status := range ss {
				if status == "dead" {
					reIps = append(reIps, strings.Split(ip, ":")[0])
				}
			}
			if len(reIps) > 0 {
				AgentDestDir := beego.AppConfig.String("AgentDestDir")
				s.StartP2pAgent(reIps, AgentDestDir)
			}
		}
	}
	beego.Info("p2p agent自动任务结束" + time.Now().Format("2006-01-02 15:04:05"))
	return nil
}
