package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"library/common"
	"models"
)

type TestApiController struct {
	beego.Controller
}

func (c *TestApiController) Get() {
	var projects []models.Project
	o := orm.NewOrm()
	o.Raw("SELECT * FROM `project`  WHERE 1=1").QueryRows(&projects)
	for _, v := range projects {
		if v.PmsProName != "" {
			var projectGits []orm.Params
			ss, err := o.Raw("SELECT id,name FROM pms.project_git WHERE `name`=?", v.PmsProName).Values(&projectGits)
			if err == nil && ss > 0 {
				v.PmsProName = common.GetString(projectGits[0]["id"])
				if v.PmsProName != "" {
					err = models.UpdateProjectById(&v)
				}
			}
		}
	}

}
