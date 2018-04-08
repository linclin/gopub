package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"library/p2p/init_sever"
	"models"
	"os"
	"os/signal"
	_ "routers"
	"syscall"
	"tasks"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
	_ "github.com/go-sql-driver/mysql"
)

func initArgs() {
	args := os.Args
	for _, v := range args {
		if v == "-syncdb" {
			models.Syncdb()
			os.Exit(0)
		}
		if v == "-docker" {
			beego.BConfig.RunMode = "docker"
			models.Syncdb()
		}
	}
}

func init() {
	//初始化数据库
	initArgs()
	//连接MySQL
	dbUser := beego.AppConfig.String("mysqluser")
	dbPass := beego.AppConfig.String("mysqlpass")
	dbHost := beego.AppConfig.String("mysqlhost")
	dbPort := beego.AppConfig.String("mysqlport")
	dbName := beego.AppConfig.String("mysqldb")
	maxIdleConn, _ := beego.AppConfig.Int("mysql_max_idle_conn")
	maxOpenConn, _ := beego.AppConfig.Int("mysql_max_open_conn")
	dbLink := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", dbUser, dbPass, dbHost, dbPort, dbName) + "&loc=Asia%2FShanghai"
	//utils.Display("dbLink", dbLink)
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", dbLink, maxIdleConn, maxOpenConn)

	if beego.BConfig.RunMode == "dev" {
		orm.Debug = true
	}
	//设置日志
	fn := "logs/run.log"
	if _, err := os.Stat(fn); err != nil {
		if os.IsNotExist(err) {
			os.Create(fn)
		}
	}
	beego.SetLogger("file", `{"filename":"`+fn+`"}`)
	if beego.BConfig.RunMode == "prod" {
		beego.SetLevel(beego.LevelInformational)
	}

}

func handleSignals(c chan os.Signal) {
	switch <-c {
	case syscall.SIGINT, syscall.SIGTERM:

		beego.Info("Shutdown quickly, bye...")
	case syscall.SIGQUIT:
		beego.Info("Shutdown gracefully, bye...")
		// do graceful shutdown
	}
	os.Exit(0)
}

func main() {
	//获取全局panic
	defer func() {
		if err := recover(); err != nil {
			beego.Error("Panic error:", err)
		}
	}()
	//热启动
	graceful, _ := beego.AppConfig.Bool("Graceful")
	if graceful {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		go handleSignals(sigs)
	}
	//API自动化文档
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"

	}

	if beego.BConfig.RunMode == "prod" {
		//check_p2p_angent_status := toolbox.NewTask("check_p2p_angent_status", "0 0 0 * * 0", func() error {
		//	err := tasks.Check_p2p_angent_status()
		//	if err != nil {
		//		beego.Error("定时任务: check_p2p_angent_status 发生错误:", err.Error())
		//		return err
		//	}
		//	return nil
		//})
		//toolbox.AddTask("check_p2p_angent_status", check_p2p_angent_status)
		defer toolbox.StopTask()
	}

	toolbox.StartTask()

	//开启双向SSL认证,认证客户端证书
	VerifyClientCert, _ := beego.AppConfig.Bool("VerifyClientCert")
	if VerifyClientCert {

		pool := x509.NewCertPool()
		caCertPath := beego.AppConfig.String("CaCertFile")
		caCrt, err := ioutil.ReadFile(caCertPath)
		if err != nil {
			panic("CA File 读取错误::" + err.Error())
		}
		pool.AppendCertsFromPEM(caCrt)
		config := tls.Config{
			ClientAuth: tls.RequireAndVerifyClientCert,
			//Certificates: []tls.Certificate{cert},
			ClientCAs: pool,
		}
		config.Rand = rand.Reader
		beego.BeeApp.Server.TLSConfig = &config
	}
	beego.Info(beego.BConfig.RunMode)
	if beego.BConfig.RunMode != "docker" {
		init_sever.Start()
	}else if  beego.BConfig.RunMode != "init"   {

	}
	beego.Run()
}
