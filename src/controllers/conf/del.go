package confcontrollers

import (
	"github.com/linclin/gopub/src/controllers"
	"github.com/linclin/gopub/src/models"
)

type DelController struct {
	controllers.BaseController
}

func (c *DelController) Get() {
	projectId, _ := c.GetInt("projectId", 0)
	err := models.DeleteProject(projectId)
	if err != nil {
		c.SetJson(1, nil, err.Error())
		return
	}
	c.SetJson(0, nil, "删除成功")
	return
}
