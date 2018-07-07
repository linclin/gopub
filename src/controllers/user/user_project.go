package usercontrollers

import (
	"controllers"
	"github.com/astaxie/beego/orm"
)

type UserProjectController struct {
	controllers.BaseController
}

func (c *UserProjectController) Get() {
	userId := c.GetString("user_id")
	o := orm.NewOrm()
	var projects []orm.Params
	o.Raw("SELECT  project.id,project.`name`,project.`level` FROM `group` left JOIN project on  group.project_id=project.id WHERE `group`.user_id= ?", userId).Values(&projects)
	c.SetJson(0, projects, "")
	return

}
