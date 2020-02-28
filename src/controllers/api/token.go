package apicontrollers

import (
	"github.com/linclin/gopub/src/library/common"
	"github.com/linclin/gopub/src/models"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
)

// 访问Token管理
type TokenController struct {
	beego.Controller
}

func (c *TokenController) URLMapping() {
	c.Mapping("GetOne", c.IssueToken)
}

// @Title 生成访问Token
// @Description 生成Json Web Token
// @Param	appid	query	string	false	"appid"
// @Param	appsecret	query	string	false	"appsecret"
// @Success 200   {"access_token":"ACCESS_TOKEN","expires_in":"7200"}
// @Failure 200   {"errcode":"100","errmsg":"invalid appid"}
// @router / [get]
func (c *TokenController) IssueToken() {

	appid := c.GetString("appid", "")
	appsecret := c.GetString("appsecret", "")
	if appid == "" || appsecret == "" {
		c.Data["json"] = map[string]string{"errcode": "100", "errmsg": "appid & appsecret必须 "}
		c.ServeJSON()
		return
	}
	clientip := c.Ctx.Input.IP()
	api_system, err := models.GetApiSystemById(common.StrToInt(appid))
	if err != nil {
		c.Data["json"] = map[string]string{"errcode": "100", "errmsg": "appid不存在 "}
		c.ServeJSON()
		return
	}
	if api_system.AppSecret != appsecret {
		c.Data["json"] = map[string]string{"errcode": "100", "errmsg": "appid 或者 appsecret不匹配  "}
		c.ServeJSON()
		return
	}
	api_system_ips := strings.Split(api_system.IP, ",")
	beego.Debug(api_system_ips)
	ipin := false
	for _, ip := range api_system_ips {
		if clientip == ip {
			ipin = true
		}
	}
	if ipin == false {
		c.Data["json"] = map[string]string{"errcode": "100", "errmsg": clientip + " 请求IP不在允许范围内  "}
		c.ServeJSON()
		return
	}
	// Create a Token that will be signed with HS256.
	//nowtime := time.Now().Unix()
	exptime := time.Now().Unix() + 3600
	token := jwt.New(jwt.SigningMethodHS256)
	claims, _ := token.Claims.(jwt.MapClaims)
	claims["iss"] = appid //The issuer of the token，token 是给谁的
	//token.Claims["sub"] = beego.AppConfig.String("AppName")                // The subject of the token，token 主题
	//token.Claims["jti"] = appid + strconv.FormatInt(time.Now().Unix(), 10) //JWT ID。针对当前 token 的唯一标识
	//token.Claims["iat"] = nowtime //Issued At。 token 创建时间， Unix 时间戳格式
	claims["exp"] = exptime //Expiration Time。 token 过期时间，Unix 时间戳格式
	// The claims object allows you to store information in the actual token.
	tokenString, err := token.SignedString([]byte(beego.AppConfig.String("SecretKey")))
	// tokenString Contains the actual token you should share with your client.
	if err != nil {
		c.Data["json"] = map[string]string{"errcode": "100", "errmsg": "Token生成错误:" + err.Error()}
	} else {
		c.Data["json"] = map[string]string{"access_token": tokenString, "expires_in": strconv.FormatInt(exptime, 10)}
	}
	c.ServeJSON()
}
