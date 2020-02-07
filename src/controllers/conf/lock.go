package confcontrollers

import (
	"controllers"
	"models"
)

type LockController struct {
	controllers.BaseController
}

func (c *LockController) Get() {
	if c.User == nil || c.User.Id == 0 {
		c.SetJson(2, nil, "not login")
		return
	}
	projectId, _ := c.GetInt("projectId", 0)
	// 1为锁定 0为解锁
	act, _ := c.GetInt("act", 0)

	project, err := models.GetProjectById(projectId)
	if err != nil {
		c.SetJson(1, nil, err.Error())
	}

	if act == 1 {
		project.UserLock = int(c.User.Id)
	} else {
		project.UserLock = 0
	}

	err = models.UpdateProjectById(project)
	if err != nil {
		c.SetJson(1, nil, err.Error())
		return
	}
	c.SetJson(0, nil, "锁定成功")
	return
}
