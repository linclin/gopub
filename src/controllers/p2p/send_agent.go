package p2pcontrollers

import (
	"controllers"
	"library/components"
	"models"

	"github.com/astaxie/beego"
)

type SendAgentController struct {
	controllers.BaseController
}

func (c *SendAgentController) Get() {
	if c.Project == nil || c.Project.Id == 0 {
		c.SetJson(1, nil, "Parameter error")
		return
	}

	s := components.BaseComponents{}
	s.SetProject(c.Project)
	s.SetTask(&models.Task{Id: -3})
	agentDir := beego.AppConfig.String("AgentDir")
	AgentDestDir := beego.AppConfig.String("AgentDestDir")
	err := s.SendP2pAgent(agentDir, AgentDestDir)
	if err != nil {
		beego.Info("出错啦！")
		c.SetJson(1, nil, "p2p文件传输失败，请检查配置，或目标机器权限"+err.Error())
		return
	}
	c.SetJson(0, nil, "更新agent成功")
	return
}
