package p2pcontrollers

import (
	"encoding/json"
	"github.com/linclin/gopub/src/controllers"
	"github.com/linclin/gopub/src/library/components"
	"github.com/linclin/gopub/src/library/p2p/init_sever"
	"github.com/linclin/gopub/src/models"
	"strings"

	"github.com/astaxie/beego"
)

type AgentController struct {
	controllers.BaseController
}

func (c *AgentController) Get() {
	if c.Project == nil || c.Project.Id == 0 {
		c.SetJson(1, nil, "Parameter error")
		return
	}

	s := components.BaseComponents{}
	s.SetProject(c.Project)
	s.SetTask(&models.Task{Id: -3})
	ips := s.GetHostIps()
	ss := init_sever.P2pSvc.CheckAllClient(ips)
	reIps := []string{}
	for ip, status := range ss {
		if status == "dead" {
			reIps = append(reIps, strings.Split(ip, ":")[0])
		}
	}
	if len(reIps) > 0 && c.Project.P2p == 1 {
		AgentDestDir := beego.AppConfig.String("AgentDestDir")
		err := s.StartP2pAgent(reIps, AgentDestDir)
		if err != nil {
			c.SetJson(1, nil, "重启失败"+err.Error())
			return
		} else {
			c.SetJson(0, nil, "重启成功")
			return
		}
	} else {
		c.SetJson(0, nil, "已全部启动")
		return
	}

	return
}

func (c *AgentController) Post() {
	if c.Project == nil || c.Project.Id == 0 {
		c.SetJson(1, nil, "Parameter error")
		return
	}
	ips := []string{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &ips)
	s := components.BaseComponents{}
	s.SetProject(c.Project)
	AgentDestDir := beego.AppConfig.String("AgentDestDir")
	err := s.StartP2pAgent(ips, AgentDestDir)
	if err != nil {
		c.SetJson(1, nil, "重启失败"+err.Error())
		return
	} else {
		c.SetJson(0, nil, "重启成功")
		return
	}
}
