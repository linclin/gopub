package p2pcontrollers

import (
	"controllers"
	"library/p2p/init_sever"
)

type TaskController struct {
	controllers.BaseController
}

func (c *TaskController) Get() {
	taskId := c.GetString("taskId")
	ss, _ := init_sever.P2pSvc.QueryTaskNoHttp(taskId)
	c.SetJson(0, ss, "")
	return
}
