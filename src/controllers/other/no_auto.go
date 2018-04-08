package othercontrollers

import (
	"controllers"
	"fmt"
	"library/common"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type NoAutoController struct {
	controllers.BaseController
}

//这里是查询每天 每周 每月 未进入预发布的项目
func (c *NoAutoController) Get() {
	taskType := c.GetString("taskType")
	beego.Info(taskType)
	o := orm.NewOrm()
	sql := "SELECT project.id ,project.name  FROM `task` LEFT JOIN project ON task.project_id=project.id WHERE  project.level=2 %s group BY project.id"
	var proIds []orm.Params
	var proIds1 []orm.Params
	if taskType == "day" {
		o.Raw(fmt.Sprintf(sql, " and TO_DAYS(now()) - TO_DAYS(task.updated_at) = 0 ")).Values(&proIds)
		o.Raw(fmt.Sprintf(sql, " and TO_DAYS(now()) - TO_DAYS(task.updated_at) = 0  and task.user_id=1")).Values(&proIds)
	} else if taskType == "week" {
		o.Raw(fmt.Sprintf(sql, " and YEARWEEK(date_format(task.updated_at,'%Y-%m-%d')) = YEARWEEK(now()) ")).Values(&proIds)
		o.Raw(fmt.Sprintf(sql, " and YEARWEEK(date_format(task.updated_at,'%Y-%m-%d')) = YEARWEEK(now()) and task.user_id=1")).Values(&proIds1)
	} else {
		o.Raw(fmt.Sprintf(sql, " and date_format(task.updated_at,'%Y-%m')=date_format(now(),'%Y-%m') ")).Values(&proIds)
		o.Raw(fmt.Sprintf(sql, " and date_format(task.updated_at,'%Y-%m')=date_format(now(),'%Y-%m')  and task.user_id=1")).Values(&proIds1)
	}
	res := []orm.Params{}
	for _, proId := range proIds {
		id := common.GetInt(proId["id"])
		isIn := false
		for _, proId1 := range proIds1 {
			id1 := common.GetInt(proId1["id"])
			if id == id1 {
				isIn = true
			}
		}
		if !isIn {
			res = append(res, proId)
		}
	}
	c.SetJson(0, res, "")
	return

}
