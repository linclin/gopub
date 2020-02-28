package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/linclin/gopub/src/models"
)

type TestApiController struct {
	beego.Controller
}

func (c *TestApiController) Get() {
	var projects []models.Project
	o := orm.NewOrm()
	o.Raw("SELECT * FROM `project`  WHERE 1=1").QueryRows(&projects)
	c.Data["json"] = projects
	c.ServeJSON()

}
