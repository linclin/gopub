package p2pcontrollers

import (
	"controllers"
	"library/common"
	"library/components"
	"library/p2p/init_sever"
	"models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type CheckController struct {
	controllers.BaseController
}

type P2pinfo struct {
	Host   string
	Status string
	Pid    int
	Pname  string
}

func (c *CheckController) Get() {
	searchtype := c.GetString("type")
	projectId := c.GetString("projectId")
	beego.Info(searchtype)
	if searchtype == "0" {
		o := orm.NewOrm()
		var projects []models.Project
		var p []P2pinfo
		ss := map[string]string{}
		i, err := o.Raw("SELECT * FROM `project` WHERE `p2p` = 1 ").QueryRows(&projects)
		if i > 0 && err == nil {
			for _, project := range projects {
				s := components.BaseComponents{}
				s.SetProject(&project)
				ips := s.GetHostIps()
				proRes := init_sever.P2pSvc.CheckAllClient(ips)
				for key, value := range proRes {
					if value == "dead" {
						pa := P2pinfo{}
						if !common.InList(key, ss) {
							ss[key] = value
							pa.Host = key
							pa.Status = value
							pa.Pid = project.Id
							pa.Pname = project.Name
							p = append(p, pa)
						}

					}
				}
			}
			beego.Info(p)
			c.SetJson(0, p, "")
			return
		} else {
			c.SetJson(1, ss, "no agent")
			return
		}
	} else if projectId != "" && searchtype == "1" {
		o := orm.NewOrm()
		var projects []models.Project
		ss := map[string]string{}
		i, err := o.Raw("SELECT * FROM `project` WHERE `id` = ?   ", projectId).QueryRows(&projects)
		if i > 0 && err == nil {
			for _, project := range projects {
				s := components.BaseComponents{}
				s.SetProject(&project)
				ips := s.GetHostIps()
				proRes := init_sever.P2pSvc.CheckAllClient(ips)
				for key, value := range proRes {
					if !common.InList(key, ss) {
						ss[key] = value
					}
				}
			}
			c.SetJson(0, ss, "")
			return
		} else {
			c.SetJson(1, ss, "no agent")
			return
		}
	}
	return
}
