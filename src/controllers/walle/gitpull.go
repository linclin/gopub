package wallecontrollers

import (
	"controllers"
	"library/components"
	"models"
)

type GitpullController struct {
	controllers.BaseController
}

func (c *GitpullController) Get() {
	if c.Project == nil || c.Project.Id == 0 {
		c.SetJson(1, nil, "Parameter error")
		return
	}
	s := components.BaseComponents{}
	s.SetProject(c.Project)
	s.SetTask(&models.Task{Id: -3})
	err := s.GetGitPull()
	if err != nil {
		c.SetJson(1, nil, "拉取错误—"+err.Error())
		return
	} else {
		c.SetJson(0, nil, "")
		return
	}

}
