package confcontrollers

import (
	"controllers"
	"github.com/astaxie/beego"
	"library/common"

	"library/jumpserver"
	"strings"
)

type GroupInfoController struct {
	controllers.BaseController
}

func (c *GroupInfoController) Get() {
	if c.User == nil || c.User.Id == 0 {
		c.SetJson(2, nil, "not login")
		return
	}
	groupid := c.GetString("hostgroup")
	if groupid == "" {
		c.SetJson(1, nil, "params")
	}
	aGroupid := strings.Split(groupid, " ")
	if len(aGroupid) < 1 {
		c.SetJson(1, nil, "params array")
	}

	mGroupid2true := make(map[string]bool)
	var rsIps []string
	for _, gid := range aGroupid {
		aIp, _ := jumpserver.GetIpsByGroupid(string(gid))
		beego.Info(aIp)
		mGroupid2true[string(gid)] = true
		if len(aIp) > 0 {
			for ip, _ := range aIp {
				rsIps = append(rsIps, ip)
			}
		}
	}
	rsIps = common.ArrayUnique(rsIps)

	rsId2Groupname := make(map[string]string)
	group2id, _ := jumpserver.GetGroups()
	if len(group2id) > 0 {
		for group_id, groupname := range group2id {
			_, ok := mGroupid2true[group_id]
			if ok == true {
				rsId2Groupname[group_id] = groupname
			}

		}
	}

	rs := make(map[string]interface{})
	rs["id2groupname"] = rsId2Groupname
	rs["ips"] = rsIps

	c.SetJson(0, rs, "")
	return
}
