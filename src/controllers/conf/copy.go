package confcontrollers

import (
	"github.com/linclin/gopub/src/controllers"
	"github.com/linclin/gopub/src/models"
)

type CopyController struct {
	controllers.BaseController
}

func (c *CopyController) Get() {
	if c.User == nil || c.User.Id == 0 {
		c.SetJson(2, nil, "not login")
		return
	}
	if c.Project == nil || c.Project.Id == 0 {
		c.SetJson(1, nil, "Parameter error")
		return
	}
	c.Project.Name = c.Project.Name + " - copy"
	c.Project.Id = 0
	c.Project.UserId = uint(c.User.Id)
	_, err := models.AddProject(c.Project)
	if err != nil {
		c.SetJson(1, nil, "复制失败")
	}
	c.SetJson(0, c.Project, "")
	return

}
