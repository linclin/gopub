package wallecontrollers

import (
	"controllers"
	"library/components"
	"models"
)

type CommitController struct {
	controllers.BaseController
}

func (c *CommitController) Get() {
	if c.Project == nil || c.Project.Id == 0 {
		c.SetJson(1, nil, "Parameter error")
		return
	}
	branch := c.GetString("branch")
	s := components.BaseComponents{}
	s.SetProject(c.Project)
	s.SetTask(&models.Task{})
	g := components.BaseGit{}
	g.SetBaseComponents(s)
	res, err := g.GetCommitList(branch, 25)
	if err != nil {
		c.SetJson(1, nil, "获取Commit错误—"+err.Error())
		return
	} else {
		c.SetJson(0, res, "")
		return
	}

}
