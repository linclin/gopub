package taskcontrollers

import (
	"controllers"
	"github.com/astaxie/beego/orm"
	"models"
)

type LastTaskController struct {
	controllers.BaseController
}

func (c *LastTaskController) Get() {
	if c.Project == nil || c.Project.Id == 0 {
		c.SetJson(1, nil, "Parameter error")
		return
	}
	o := orm.NewOrm()
	var task models.Task
	o.Raw("SELECT * FROM task where project_id = ? AND status=3 order by task.id DESC LIMIT 1", c.Project).QueryRow(&task)
	c.SetJson(0, task, "")
	return

}
