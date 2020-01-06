package confcontrollers

import (
	"controllers"
	"github.com/astaxie/beego/orm"
	"library/common"
)

type MyListController struct {
	controllers.BaseController
}

func (c *MyListController) Get() {
	page, _ := c.GetInt("page", 0)
	start := 0
	length, _ := c.GetInt("length", 200000)
	if page > 0 {
		start = (page - 1) * length
	}
	selectInfo := c.GetString("select_info")
	where := ""
	if selectInfo != "" {
		where = "  and(`name` LIKE '%" + selectInfo + "%' )"
	}
	var projects []orm.Params
	o := orm.NewOrm()
	if c.User.Role == 10 {
		where = where + "and  `level`= 2  "
	} else if c.User.Role == 20 {
		where = where + "and  id in (SELECT project_id FROM `group` WHERE `group`.user_id=" + common.GetString(c.User.Id) + " )  "

	}
	o.Raw("SELECT *, (SELECT realname FROM `user` WHERE `user`.id=project.user_id LIMIT 1) as realname,(SELECT realname FROM `user` WHERE `user`.id=project.user_lock LIMIT 1) as lockuser FROM `project`  WHERE 1=1 "+where+" ORDER BY id LIMIT ?,?", start, length).Values(&projects)
	var count []orm.Params
	total := 0
	o.Raw("SELECT count(id) FROM `project` WHERE 1=1 " + where).Values(&count)
	if len(count) > 0 {
		total = common.GetInt(count[0]["count(id)"])
	}
	c.SetJson(0, map[string]interface{}{"total": total, "currentPage": page, "table_data": projects}, "")

	return

}
