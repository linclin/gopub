package controllers

import (
	"library/common"
	"runtime"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"models"
	"strings"
)

//基类
type BaseController struct {
	beego.Controller
	Project *models.Project
	Task    *models.Task
	User    *models.User
}

// Prepare implemented Prepare method for baseRouter.
func (c *BaseController) Prepare() {
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")                                                      //允许访问源
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")                               //允许post访问
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "X-Requested-With,UserToken,Content-Type,Authorization") //header的类型
	if c.Ctx.Input.IsOptions() {
		return
	}

	//获取panic
	defer func() {
		if panic_err := recover(); panic_err != nil {
			var buf []byte = make([]byte, 1024)
			runtimec := runtime.Stack(buf, false)
			beego.Error("控制器错误:", panic_err, string(buf[0:runtimec]))
		}
	}()
	taskId := ""
	if c.Ctx.Input.Param(":taskId") != "" {
		taskId = c.Ctx.Input.Param(":taskId")
	} else if c.GetString("taskId") != "" {
		taskId = c.GetString("taskId")
	}
	if taskId != "" {
		c.Task, _ = models.GetTaskById(common.GetInt(taskId))
	}
	projectId := ""
	if c.Ctx.Input.Param(":projectId") != "" {
		projectId = c.Ctx.Input.Param(":projectId")
	} else if c.GetString("projectId") != "" {
		projectId = c.GetString("projectId")
	}
	if projectId != "" {
		c.Project, _ = models.GetProjectById(common.GetInt(projectId))
	}
	token := ""
	if ah := c.Ctx.Input.Header("Authorization"); ah != "" {
		if len(ah) > 5 && strings.ToUpper(ah[0:5]) == "TOKEN" {
			token = ah[6:]
			if token != "" {
				var users []models.User
				o := orm.NewOrm()
				s, err := o.Raw("SELECT * FROM `user` WHERE auth_key= ?", token).QueryRows(&users)
				if s > 0 && err == nil {
					c.User = &(users[0])
				}
			}
		}
	}
}
func (this *BaseController) SetJson(code int, data interface{}, Msg string) {
	if code == 0 {
		if Msg == "" {
			Msg = "sucess"
		}
		this.Data["json"] = map[string]interface{}{"code": code, "msg": Msg, "data": data}
		this.ServeJSON()
	} else {
		this.Data["json"] = map[string]interface{}{"code": code, "msg": Msg, "data": data}
		this.ServeJSON()
	}

}
func (this *BaseController) AllowCross() {
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")                                                      //允许访问源
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")                               //允许post访问
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "X-Requested-With,UserToken,Content-Type,Authorization") //header的类型
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Max-Age", "1728000")
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	this.Ctx.ResponseWriter.Header().Set("content-type", "application/json") //返回数据格式是json
}
func (c *BaseController) Options() {
	c.AllowCross() //允许跨域
	c.Data["json"] = map[string]interface{}{"status": 200, "message": "ok", "moreinfo": ""}
	c.ServeJSON()
}
