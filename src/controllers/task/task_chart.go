package taskcontrollers

import (
	"controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	"library/common"
	"library/components"
	"models"
	"time"
)

type TaskChartController struct {
	controllers.BaseController
}

var bm, _ = cache.NewCache("memory", `{"interval":3600}`)

func (c *TaskChartController) Get() {
	taskType := c.GetString("taskType")
	o := orm.NewOrm()
	if taskType == "day" {
		var count []orm.Params
		o.Raw("SELECT project.`level`,count(task.id) as task_count  FROM `task` LEFT JOIN project ON task.project_id = project.id WHERE TO_DAYS(now()) - TO_DAYS(task.updated_at) = 0 GROUP BY project. LEVEL").Values(&count)
		for _, c := range count {
			c["name"] = GetProjectLevel(common.GetInt(c["level"]))
		}
		c.SetJson(0, count, "")
		return
	} else if taskType == "week" {
		var count []orm.Params
		o.Raw("SELECT project.`level`,count(task.id) as task_count  FROM `task` LEFT JOIN project ON task.project_id = project.id WHERE YEARWEEK(date_format(task.updated_at,'%Y-%m-%d')) = YEARWEEK(now()) GROUP BY project. LEVEL").Values(&count)
		for _, c := range count {
			c["name"] = GetProjectLevel(common.GetInt(c["level"]))
		}
		c.SetJson(0, count, "")
		return
	} else if taskType == "month" {
		var count []orm.Params
		o.Raw("SELECT project.`level`,count(task.id) as task_count  FROM `task` LEFT JOIN project ON task.project_id = project.id WHERE date_format(task.updated_at,'%Y-%m')=date_format(now(),'%Y-%m') GROUP BY project. LEVEL").Values(&count)
		for _, c := range count {
			c["name"] = GetProjectLevel(common.GetInt(c["level"]))
		}
		beego.Info(count)
		c.SetJson(0, count, "")
		return
	} else if taskType == "dayBypro" {
		var count []orm.Params
		o.Raw("SELECT project.`name`,count(task.id) as task_count,project.`level` FROM `task` LEFT JOIN project ON task.project_id = project.id WHERE TO_DAYS(now()) - TO_DAYS(task.updated_at) = 0 and task.status=3 GROUP BY project.id").Values(&count)
		for _, c := range count {
			c["name"] = common.GetString(c["name"]) + "-" + GetProjectLevel(common.GetInt(c["level"]))
		}
		c.SetJson(0, count, "")
		return
	} else if taskType == "weekBypro" {
		var count []orm.Params
		o.Raw("SELECT project.`name`,count(task.id) as task_count,project.`level` FROM `task` LEFT JOIN project ON task.project_id = project.id WHERE YEARWEEK(date_format(task.updated_at,'%Y-%m-%d')) = YEARWEEK(now()) and task.status=3 GROUP BY project.id").Values(&count)
		for _, c := range count {
			c["name"] = common.GetString(c["name"]) + "-" + GetProjectLevel(common.GetInt(c["level"]))
		}
		c.SetJson(0, count, "")
		return
	} else if taskType == "monthBypro" {
		var count []orm.Params
		o.Raw("SELECT project.`name`,count(task.id) as task_count,project.`level` FROM `task` LEFT JOIN project ON task.project_id = project.id WHERE date_format(task.updated_at,'%Y-%m')=date_format(now(),'%Y-%m') and task.status=3 GROUP BY project.id").Values(&count)
		for _, c := range count {
			c["name"] = common.GetString(c["name"]) + "-" + GetProjectLevel(common.GetInt(c["level"]))
		}
		c.SetJson(0, count, "")
		return
	} else if taskType == "total" {
		totalJson := map[string]interface{}{}
		var totalmem []orm.Params
		var totalproject []orm.Params
		var totalpub []orm.Params
		var totalpubsuccess []orm.Params
		num, err := o.Raw("SELECT count(id) as `totalmen` FROM `user`").Values(&totalmem)
		if num > 0 && err == nil {
			totalJson["totalmen"] = common.GetInt(totalmem[0]["totalmen"])
		}
		num, err = o.Raw("SELECT count(DISTINCT name) as `totalproject` from `project`").Values(&totalproject)
		if num > 0 && err == nil {
			totalJson["totalproject"] = common.GetInt(totalproject[0]["totalproject"])
		}
		num, err = o.Raw("SELECT count(id) as `totalpub` from `task`").Values(&totalpub)
		if num > 0 && err == nil {
			totalJson["totalpub"] = common.GetInt(totalpub[0]["totalpub"])
		}
		num, err = o.Raw("SELECT count(id) as `totalpubsuccess` from `task`where status = 3").Values(&totalpubsuccess)
		if num > 0 && err == nil {
			totalJson["totalpubsuccess"] = common.GetInt(totalpubsuccess[0]["totalpubsuccess"])
		}
		if bm.IsExist("hostsum") == false {
			totalJson["hostsum"] = GetHostNum()
		} else {
			totalJson["hostsum"] = bm.Get("hostsum")
		}
		c.SetJson(0, totalJson, "")
		return
	}
	c.SetJson(1, nil, "未传参数")
	return

}
func GetProjectLevel(level int) string {
	switch level {
	case 1:
		return "测试环境"
		break
	case 2:
		return "预发布环境"
		break
	case 3:
		return "生产环境"
		break
	}
	return "删除项目"
}

func GetHostNum() int {
	o := orm.NewOrm()
	var projects []models.Project
	i, err := o.Raw("SELECT * FROM `project`").QueryRows(&projects)
	finalres := []string{}
	if i > 0 && err == nil {
		for _, project := range projects {
			s := components.BaseComponents{}
			s.SetProject(&project)
			ips := s.GetHosts()
			for _, ip := range ips {
				if !common.InList(string(ip), finalres) {
					finalres = append(finalres, string(ip))
				}
			}
		}
	}
	bm.Put("hostsum", len(finalres), 1*time.Hour)
	return (len(finalres))
}
