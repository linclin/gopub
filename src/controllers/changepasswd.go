package controllers

import (
	"github.com/astaxie/beego"

	"encoding/json"
	"github.com/astaxie/beego/orm"
	"github.com/linclin/gopub/src/library/common"
	"github.com/linclin/gopub/src/models"
	"golang.org/x/crypto/bcrypt"
)

type ChangePasswdController struct {
	BaseController
}

func (c *ChangePasswdController) Post() {
	//哈希校验成功后 更新 auth_key
	if c.User == nil || c.User.Id == 0 {
		c.SetJson(2, nil, "not login")
		return
	}

	beego.Info(string(c.Ctx.Input.RequestBody))

	postData := map[string]string{"newpassword": "", "repeat_newpassword": ""}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &postData)
	if err != nil {
		c.SetJson(1, nil, "数据格式错误")
		return
	}

	uid := postData["uid"]
	if common.GetString(c.User.Id) != uid && c.User.Role != 1 {
		c.SetJson(1, nil, "403")
		return
	}
	newPassword := postData["newpassword"]
	repeatNewpassword := postData["repeat_newpassword"]
	if newPassword == "" || repeatNewpassword == "" {
		c.SetJson(1, nil, "请输入密码")
		return
	}
	var user models.User
	o := orm.NewOrm()
	err = o.Raw("SELECT * FROM `user` WHERE id= ?", uid).QueryRow(&user)
	beego.Info(err)
	//验证旧密码

	if newPassword == repeatNewpassword {
		password := []byte(newPassword)
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		user.PasswordHash = string(hashedPassword)
		models.UpdateUserById(&user)
		c.Data["json"] = map[string]interface{}{"code": 0, "msg": "sucess"}
		c.ServeJSON()
		return
	} else {
		c.SetJson(1, nil, "两次密码输入不一致，请重新输入")
		return
	}
}
