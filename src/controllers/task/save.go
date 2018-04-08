package taskcontrollers

import (
	"controllers"
	"encoding/json"
	"github.com/astaxie/beego"
	"models"
	"time"
)

type SaveController struct {
	controllers.BaseController
}

func (c *SaveController) Post() {
	//projectId,_:=c.GetInt("projectId",0)
	if c.User == nil || c.User.Id == 0 {
		c.SetJson(2, nil, "not login")
		return
	}
	beego.Info(string(c.Ctx.Input.RequestBody))
	var task models.Task
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &task)
	if err != nil {
		c.SetJson(1, nil, "数据库更新错误"+err.Error())
	}
	if task.Id != 0 {
		err = models.UpdateTaskById(&task)
	} else {
		if task.Hosts == "" {
			ss, err := models.GetProjectById(task.ProjectId)
			if err == nil {
				task.Hosts = ss.Hosts
			}
		}
		task.UserId = uint(c.User.Id)
		task.CreatedAt = time.Now()
		task.UpdatedAt = time.Now()
		task.EnableRollback = 1
		_, err = models.AddTask(&task)
	}
	if err != nil {
		c.SetJson(1, nil, "数据库更新错误")
	}
	c.SetJson(0, task, "修改成功")

	return
}
