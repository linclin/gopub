package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/bcrypt"
	"library/common"
	"models"
	"regexp"
	"strings"
	"time"
)

type RegisterController struct {
	BaseController
}

//邮箱正则
func IsEmail(str ...string) bool {
	var b bool
	for _, s := range str {
		b, _ = regexp.MatchString("^([a-z0-9_\\.-]+)@([\\da-z\\.-]+)\\.([a-z\\.]{2,6})$", s)
		if false == b {
			return b
		}
	}
	return b
}
func (c *RegisterController) Post() {

	beego.Info(string(c.Ctx.Input.RequestBody))
	registerData := map[string]interface{}{"user_password": "", "user_name": "", "Role": 1}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &registerData)
	if err != nil {
		c.SetJson(1, nil, "数据格式错误")
		return
	}
	registerUsername := common.GetString(registerData["register_username"])
	registerRealname := common.GetString(registerData["register_realname"])
	registerEmail := common.GetString(registerData["register_email"])

	registerRole := common.GetInt(registerData["Role"])
	//格式判断
	realnmae := strings.Split(registerRealname, ".")
	if len(realnmae) != 2 {
		c.SetJson(1, nil, "花名.实名输入有误")
		return
	}

	iseamil := IsEmail(registerEmail)
	if iseamil == false {
		c.SetJson(1, nil, "邮箱输入有误")
		return
	}

	var user models.User
	o := orm.NewOrm()
	//先判断存在用户否
	err = o.Raw("SELECT * FROM `user` WHERE username= ?", registerUsername).QueryRow(&user)
	beego.Info(user)
	if err == nil {
		c.SetJson(1, nil, "用户已存在，请更换账户名")
		return
	} else { //不存在，存库
		var newuser models.User
		userAuth := common.Md5String(registerUsername + common.GetString(time.Now().Unix()))
		password := []byte("123456")
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}

		newuser.Username = registerUsername
		newuser.PasswordHash = string(hashedPassword)
		newuser.IsEmailVerified = 1
		newuser.Avatar = "default.jpg"
		newuser.Role = int16(registerRole)
		newuser.Status = 10
		newuser.CreatedAt = time.Now()
		newuser.UpdatedAt = time.Now()
		newuser.AuthKey = userAuth
		newuser.Email = registerEmail
		newuser.Realname = registerRealname

		newid, err := models.AddUser(&newuser)
		if newuser.Role == 20 {
			pro_ids := common.GetString(registerData["pro_ids"])
			pro_idArr := strings.Split(pro_ids, ",")
			for _, pro_id := range pro_idArr {
				o.Raw("INSERT INTO `group`(`project_id`, `user_id`) VALUES (?, ?)", pro_id, newid).Exec()
			}
		}

		if err != nil {
			c.SetJson(1, nil, "数据库存储错误")
			return
		} else {
			c.SetJson(0, nil, "success")
			return
		}

	}
}
