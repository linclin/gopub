package wallecontrollers

import (
	"github.com/linclin/gopub/src/controllers"
	"github.com/linclin/gopub/src/library/components"
	"github.com/linclin/gopub/src/models"
)

type GetMd5Controller struct {
	controllers.BaseController
}

func (c *GetMd5Controller) Get() {
	if c.Project == nil || c.Project.Id == 0 {
		c.SetJson(1, nil, "Parameter error")
		return
	}
	url := c.GetString("url")
	s := components.BaseComponents{}
	s.SetProject(c.Project)
	s.SetTask(&models.Task{})
	f := components.BaseFile{}
	f.SetBaseComponents(s)
	err := f.UpdateRepo(url, "")
	if err != nil {
		c.SetJson(1, nil, "获取md5错误—"+err.Error())
		return
	}
	res, err := f.CheckFiles(url, "")
	if err != nil {
		c.SetJson(1, nil, "获取md5错误—"+err.Error())
		return
	} else {
		c.SetJson(0, res, "")
		return
	}
}
