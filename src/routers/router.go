package routers

import (
	"controllers"
	"controllers/api"
	"controllers/conf"
	"controllers/other"
	"controllers/p2p"
	"controllers/record"
	"controllers/task"
	"controllers/walle"
	"github.com/astaxie/beego"
	"controllers/user"
	"github.com/astaxie/beego/plugins/cors"
	"time"
)

func init() {

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin","UserToken", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		MaxAge:           5 * time.Minute,
	}))

	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/logout", &controllers.LogoutController{})
	beego.Router("/loginbydocke", &controllers.LoginByDockerController{})
	beego.Router("/changePasswd", &controllers.ChangePasswdController{})
	beego.Router("/register", &controllers.RegisterController{})

	beego.Router("/api/get/walle/detection", &wallecontrollers.DetectionController{})
	beego.Router("/api/get/walle/detectionssh", &wallecontrollers.DetectionsshController{})
	beego.Router("/api/get/walle/release", &wallecontrollers.ReleaseController{})
	beego.Router("/api/get/walle/md5", &wallecontrollers.GetMd5Controller{})
	beego.Router("/api/get/walle/flush", &wallecontrollers.FlushController{})

	beego.Router("/api/get/conf/list", &confcontrollers.ListController{})
	beego.Router("/api/get/conf/mylist", &confcontrollers.MyListController{})
	beego.Router("/api/get/conf/get", &confcontrollers.ConfController{})
	beego.Router("/api/post/conf/save", &confcontrollers.SaveController{})
	beego.Router("/api/get/conf/del", &confcontrollers.DelController{})
	beego.Router("/api/get/conf/copy", &confcontrollers.CopyController{})

	beego.Router("/api/get/git/branch", &wallecontrollers.BranchController{})
	beego.Router("/api/get/git/commit", &wallecontrollers.CommitController{})
	beego.Router("/api/get/git/gitpull", &wallecontrollers.GitpullController{})
	beego.Router("/api/get/git/gitlog", &wallecontrollers.GitlogController{})

	beego.Router("/api/get/jenkins/commit", &wallecontrollers.JenkinsController{})

	beego.Router("/api/get/task/list", &taskcontrollers.ListController{})
	beego.Router("/api/get/task/chart", &taskcontrollers.TaskChartController{})
	beego.Router("/api/post/task/save", &taskcontrollers.SaveController{})
	beego.Router("/api/get/task/get", &taskcontrollers.TaskController{})
	beego.Router("/api/get/task/last", &taskcontrollers.LastTaskController{})
	beego.Router("/api/get/task/rollback", &taskcontrollers.RollBackController{})
	beego.Router("/api/get/task/del", &taskcontrollers.DelController{})

	beego.Router("/api/get/p2p/task", &p2pcontrollers.TaskController{})
	beego.Router("/api/get/p2p/check", &p2pcontrollers.CheckController{})
	beego.Router("/api/post/p2p/agent", &p2pcontrollers.AgentController{})
	beego.Router("/api/get/p2p/send", &p2pcontrollers.SendAgentController{})

	beego.Router("/api/get/record/list", &recordcontrollers.ListController{})

	beego.Router("/api/get/other/noauto", &othercontrollers.NoAutoController{})
	beego.Router("/api/get/test/api", &controllers.TestApiController{})
	beego.Router("/api/get/user/project", &usercontrollers.UserProjectController{})
	beego.Router("/api/get/user", &usercontrollers.UserController{})
	beego.Router("/", &controllers.MainController{})
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/token",
			beego.NSInclude(
				&apicontrollers.TokenController{},
			),
		),
		beego.NSNamespace("/task",
			beego.NSInclude(
				&apicontrollers.TaskController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
