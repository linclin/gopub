package confcontrollers

import (
	"controllers"

	"encoding/json"
	"github.com/astaxie/beego"
	"models"
)

type ConfController struct {
	controllers.BaseController
}

func (c *ConfController) Get() {
	projectId, _ := c.GetInt("projectId", 0)
	project, _ := models.GetProjectById(projectId)
	c.SetJson(0, project, "")
	return

}
func (c *ConfController) Post() {
	//projectId,_:=c.GetInt("projectId",0)
	var project models.Project
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &project)
	err = models.UpdateProjectById(&project)
	beego.Info(err)
	return
}
