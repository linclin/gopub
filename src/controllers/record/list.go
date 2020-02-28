package recordcontrollers

import (
	"github.com/astaxie/beego/orm"
	"github.com/linclin/gopub/src/controllers"
	"github.com/linclin/gopub/src/library/common"
)

type ListController struct {
	controllers.BaseController
}

func (c *ListController) Get() {
	taskId := c.GetString("taskId")
	var records []orm.Params
	if common.GetInt(taskId) <= 0 {
		timeNow := c.GetString("time")
		o := orm.NewOrm()
		o.Raw("SELECT * FROM `record` where task_id=? and created_at> ? ORDER BY `id` ASC ", taskId, timeNow).Values(&records)
	} else {
		o := orm.NewOrm()
		o.Raw("SELECT * FROM `record` where task_id=? ORDER BY `id` ASC ", taskId).Values(&records)
	}
	c.SetJson(0, records, "")
	return

}
