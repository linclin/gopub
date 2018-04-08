package taskcontrollers

import (
	"controllers"
	"models"
)

type DelController struct {
	controllers.BaseController
}

func (c *DelController) Get() {
	taskId, _ := c.GetInt("taskId", 0)
	err := models.DeleteTask(taskId)
	if err != nil {
		c.SetJson(1, nil, err.Error())
		return
	}
	c.SetJson(0, nil, "删除成功")
	return
}
