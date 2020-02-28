package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/linclin/gopub/src/library/common"
	"github.com/linclin/gopub/src/models"
	"net/url"
	"time"
)

type LoginByDockerController struct {
	BaseController
}

func (c *LoginByDockerController) Get() {
	if beego.BConfig.RunMode != "docker" {
		c.StopRun()
		return
	}
	jumpUrl := c.GetString("jumpurl")
	var user models.User
	o := orm.NewOrm()
	err := o.Raw("SELECT * FROM `user` WHERE username= ?", "admin").QueryRow(&user)

	if err == nil && user.AuthKey == "" {
		userAuth := common.Md5String(user.Username + common.GetString(time.Now().Unix()))
		user.AuthKey = userAuth
		models.UpdateUserById(&user)
	}
	resUserInfo := map[string]interface{}{"user": user, "login": true}
	userInfoJson, err := json.Marshal(resUserInfo)
	c.Ctx.SetCookie("gopub_userinfo", url.QueryEscape(string(userInfoJson)), 3600*24*2, "/")
	c.Redirect(jumpUrl, 302)
	c.StopRun()
}
