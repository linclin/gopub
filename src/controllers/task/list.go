package taskcontrollers

import (
	"controllers"
	"github.com/astaxie/beego/orm"
	"library/common"
)

type ListController struct {
	controllers.BaseController
}

func (c *ListController) Get() {
	page, _ := c.GetInt("page", 0)
	start := 0
	length, _ := c.GetInt("length", 15)
	if page > 0 {
		start = (page - 1) * length
	}
	selectInfo := c.GetString("select_info")
	where := ""
	if selectInfo != "" {
		where = "  and( project.`name` LIKE '%" + selectInfo + "%' or `user`.realname LIKE '%" + selectInfo + "%'  or task.title LIKE '%" + selectInfo + "%'  )"
	}
	myUserId, _ := c.GetInt("my", 0)
	if myUserId != 0 {
		where = where + "  and task.user_id=" + common.GetString(myUserId)
	}
	var projects []orm.Params
	o := orm.NewOrm()

	o.Raw("SELECT task.id,project.name,project.name,project.level,`user`.realname,task.title,task.action,task.link_id,task.is_run,task.enable_rollback,task.updated_at,task.branch,task.commit_id,task.pms_uwork_id,task.pms_batch_id,task.`status` FROM `task` LEFT JOIN project on task.project_id=project.id   LEFT JOIN `user` on task.user_id=user.id where 1=1 "+where+" order by task.id DESC  LIMIT ? ,?", start, length).Values(&projects)
	var count []orm.Params
	total := 0
	o.Raw("SELECT count(task.id) FROM `task` LEFT JOIN project on task.project_id=project.id   LEFT JOIN `user` on task.user_id=user.id where 1=1 " + where).Values(&count)
	if len(count) > 0 {
		total = common.GetInt(count[0]["count(task.id)"])
	}
	for _, project := range projects {
		project["status"] = GetTaskStatus(common.GetInt(project["status"]))
		if common.GetInt(project["is_run"]) != 0 && common.GetString(project["status"]) != "上线完成" {
			project["status"] = "上线中"
		}
		if common.GetInt(project["level"]) == 3 {
			project["name"] = common.GetString(project["name"]) + "-线上环境"
		}
		if common.GetInt(project["level"]) == 2 {
			project["name"] = common.GetString(project["name"]) + "-预发布环境"
		}

	}
	c.SetJson(0, map[string]interface{}{"total": total, "currentPage": page, "table_data": projects}, "")

	return

}
func GetTaskStatus(status int) string {
	switch status {
	case 0:
		return "新建提交"
		break
	case 1:
		return "新建提交"
		break
	case 2:
		return "审核拒绝"
		break
	case 3:
		return "上线完成"
		break
	case 4:
		return "上线失败"
		break
	}
	return ""
}
