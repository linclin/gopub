package controllers

import (
	"github.com/astaxie/beego"

	"encoding/json"
	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/bcrypt"
	"library/common"
	"models"
	"time"
)

type LoginController struct {
	BaseController
}

func (c *LoginController) Post() {
	//哈希校验成功后 更新 auth_key
	beego.Info(string(c.Ctx.Input.RequestBody))
	postData := map[string]string{"user_password": "", "user_name": ""}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &postData)
	if err != nil {
		c.SetJson(1, nil, "数据格式错误")
		return
	}
	password := postData["user_password"]
	userName := postData["user_name"]
	if userName == "" || password == "" {
		c.SetJson(1, nil, "用户名或密码不存在")
		return
	}
	var user models.User
	o := orm.NewOrm()
	err = o.Raw("SELECT * FROM `user` WHERE username= ?", userName).QueryRow(&user)
	beego.Info(user)
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		c.SetJson(1, nil, "用户名或密码错误")
		return
	} else {
		if user.AuthKey == "" {
			userAuth := common.Md5String(user.Username + common.GetString(time.Now().Unix()))
			user.AuthKey = userAuth
			models.UpdateUserById(&user)
		}
		user.PasswordHash=""
		c.SetJson(0, user, "")
		return
	}
}
