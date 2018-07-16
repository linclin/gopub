package taskcontrollers

import (
	"controllers"
	"encoding/json"
	"github.com/astaxie/beego"
	"library/common"
	"library/components"
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
		task.UserId = uint(c.User.Id)
		task.CreatedAt = time.Now()
		task.UpdatedAt = time.Now()
		task.EnableRollback = 1
		if task.Hosts == "" {
			ss, err := models.GetProjectById(task.ProjectId)
			if err == nil {
				task.Hosts = ss.Hosts
			}
			if ss.IsGroup == 1 {
				s := components.BaseComponents{}
				s.SetProject(ss)
				mapHost := s.GetGroupHost()
				for k, v := range mapHost {
					task1 := task
					task1.Hosts = v
					task1.Title = task1.Title + "第" + common.GetString(k) + "批"
					models.AddTask(&task1)
				}
				c.SetJson(0, task, "修改成功")
				return
			}
		}
		_, err = models.AddTask(&task)
	}
	if err != nil {
		c.SetJson(1, nil, "数据库更新错误")
	}
	c.SetJson(0, task, "修改成功")

	return
}
