package wallecontrollers

import (
	"github.com/linclin/gopub/src/controllers"
	"github.com/linclin/gopub/src/library/components"
	"github.com/linclin/gopub/src/library/p2p/init_sever"
	"github.com/linclin/gopub/src/models"
	"os"
	"time"
)

type DetectionController struct {
	controllers.BaseController
}

func (c *DetectionController) Get() {
	if c.Project == nil || c.Project.Id == 0 {
		c.SetJson(1, nil, "Parameter error")
		return
	}
	s := components.BaseComponents{}
	s.SetProject(c.Project)
	s.SetTask(&models.Task{Id: -1})
	codeBaseDir := s.GetDeployFromDir()
	//1:本地文件权限加成
	if _, err := os.Stat(codeBaseDir); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(codeBaseDir, os.ModePerm)
		}

	} else {
		_, err := os.Create(codeBaseDir + "/test.log")
		if err != nil {
			c.SetJson(1, nil, "本地文件权限错误"+err.Error())
			return
		} else {
			os.Remove(codeBaseDir + "/test.log")
		}
	}

	if c.Project.RepoType == "git" {
		//2:git信任是否添加
		g := components.BaseGit{}
		g.SetBaseComponents(s)

		err := g.UpdateRepo("", "")
		if err != nil {
			//清空后再试一次 要是不行在退出
			err := s.RemoveLocalProjectWorkspace()
			if err != nil {
				c.SetJson(1, nil, "清空目录错误"+err.Error())
				return
			}
			err = g.UpdateRepo("", "")
			if err != nil {
				c.SetJson(1, nil, "git拉取代码错误"+err.Error())
				return
			}
		}
	}
	// 3.权限与免密码登录检测
	err := s.TestSsh()
	if err != nil {
		c.SetJson(1, nil, "ssh目标机器错误"+err.Error())
		return
	}
	// 4.检测用户是否具有目标机release目录读写权限
	err = s.TestReleaseDir()
	if err != nil {
		c.SetJson(1, nil, "用户不具有目标机release目录读写权限"+err.Error())
		return
	}
	//5推送p2p客户端并启动服务
	if c.Project.P2p == 1 {
		//这里做alive检测
		ips := s.GetHostIps()
		start := time.Now()
		createdAt := int(start.Unix())
		rid := s.SaveRecord("chick p2p agent")
		ss := init_sever.P2pSvc.CheckAllClient(ips)
		for _, status := range ss {
			if status == "dead" {
				s.SaveRecordRes(rid, 0, createdAt, 0, ss)
				c.SetJson(1, nil, "p2p agent 未启动")
				return
			}
		}
		s.SaveRecordRes(rid, 0, createdAt, 1, ss)
	}
	c.SetJson(0, nil, "")
	c.ServeJSON()

}
