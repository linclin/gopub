package apicontrollers

import (
	"library/common"
	"runtime"

	"github.com/astaxie/beego"

	"github.com/dgrijalva/jwt-go"
)
//TODO 另一版的api验证 废弃
//基类
type BaseApiController struct {
	beego.Controller
}

var AppId int

// Prepare implemented Prepare method for baseRouter.
func (c *BaseApiController) Prepare() {
	//获取panic
	defer func() {
		if panic_err := recover(); panic_err != nil {
			var buf []byte = make([]byte, 1024)
			runtimec := runtime.Stack(buf, false)
			beego.Error("控制器错误:", panic_err, string(buf[0:runtimec]))

			//c.Data["json"] = map[string]string{
			//	"Error":  string(buf[0:runtimec]),
			//	"Result": common.GetString(panic_err),
			//}
			//c.ServeJSON()
			//c.StopRun()
		}
	}()

	//正式环境则需要验证Token并且需要使用HTTPS访问
	//if beego.BConfig.RunMode == "prod" {
	if c.Ctx.Input.IsSecure() == false {
		//c.Data["json"] = map[string]string{"errcode": "101", "errmsg": "请使用HTTPS请求API:" + c.Ctx.Input.Site() + ":" + beego.AppConfig.String("HttpsPort")}
		//c.ServeJSON()
		//c.StopRun()
	}
	//验证token
	tokenString := c.Ctx.Input.Header("Authorization")
	if tokenString == "" {
		c.Data["json"] = map[string]string{"errcode": "102", "errmsg": "token错误"}
		c.ServeJSON()
		c.StopRun()
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(beego.AppConfig.String("SecretKey")), nil
	})
	AppId = common.GetInt(token.Claims["iss"])
	if err != nil {
		c.Data["json"] = map[string]string{"errcode": "103", "errmsg": "token验证失败"}
		c.ServeJSON()
		c.StopRun()
	}
	//}

}
