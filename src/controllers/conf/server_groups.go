package confcontrollers

import (
	"controllers"
	"github.com/astaxie/beego"
	"library/jumpserver"
)

type ServerGroupsController struct {
	controllers.BaseController
}

func (c *ServerGroupsController) Get() {
	if c.User == nil || c.User.Id == 0 {
		c.SetJson(2, nil, "not login")
		return
	}
	enableJumpserver, _ := beego.AppConfig.Bool("enableJumpserver")
	if enableJumpserver == true {
		group2id, _ := jumpserver.GetGroups()
		c.SetJson(0, group2id, "")
	} else {
		c.SetJson(0, nil, "")
	}
	return
}
