package models

import (
	"database/sql"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"time"
)

var o orm.Ormer

func Syncdb() {
	err := createdb()
	if err != nil {
		return
	}
	Connect()
	o = orm.NewOrm()
	// 数据库别名
	name := "default"
	// drop table 后再建表
	force := true
	// 打印执行过程
	verbose := true
	// 遇到错误立即返回
	err = orm.RunSyncdb(name, force, verbose)
	if err != nil {
		fmt.Println(err)
	}
	insertUser()
	fmt.Println("database init is complete.\nPlease restart the application")

}

//数据库连接
func Connect() {
	dbUser := beego.AppConfig.String("mysqluser")
	dbPass := beego.AppConfig.String("mysqlpass")
	dbHost := beego.AppConfig.String("mysqlhost")
	dbPort := beego.AppConfig.String("mysqlport")
	dbName := beego.AppConfig.String("mysqldb")
	maxIdleConn, _ := beego.AppConfig.Int("mysql_max_idle_conn")
	maxOpenConn, _ := beego.AppConfig.Int("mysql_max_open_conn")
	dbLink := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", dbUser, dbPass, dbHost, dbPort, dbName) + "&loc=Asia%2FShanghai"
	//utils.Display("dbLink", dbLink)
	err := orm.RegisterDriver("mysql", orm.DRMySQL)
	if err != nil {
		log.Println(err)
		os.Exit(2)
		return
	}
	err = orm.RegisterDataBase("default", "mysql", dbLink, maxIdleConn, maxOpenConn)
	if err != nil {
		log.Println(err)
		os.Exit(2)
		return
	}
}

//创建数据库
func createdb() error {

	dbUser := beego.AppConfig.String("mysqluser")
	dbPass := beego.AppConfig.String("mysqlpass")
	dbHost := beego.AppConfig.String("mysqlhost")
	dbPort := beego.AppConfig.String("mysqlport")
	dbName := beego.AppConfig.String("mysqldb")

	var dsn string
	var sqlstring string

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8", dbUser, dbPass, dbHost, dbPort)
	sqlstring = fmt.Sprintf("CREATE DATABASE `%s` CHARSET utf8 COLLATE utf8_general_ci", dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println(err)
		os.Exit(2)
		//panic(err.Error())
		return err
	}
	r, err := db.Exec(sqlstring)
	if err != nil {
		log.Println(err)
		log.Println(r)
		db.Close()
		return err
	} else {
		db.Close()
		log.Println("Database ", dbName, " created")
		return nil
	}

}

func insertUser() {
	fmt.Println("insert user ...")
	u := new(User)
	u.Username = "admin"
	u.IsEmailVerified = 1
	u.AuthKey = "cJIrTa_b2Hnjn6BZkrL8PJkYto2Ael3O"
	u.PasswordHash = "$2y$13$8q0MfKpnghuqCL.3FAAjiOkA8kBFNCW.ECUlqWp1zTpMHs9e5xn6u"
	u.EmailConfirmationToken = "UpToOIawm1L8GjN6pLO4r-1oj20nLT5f_1443280741"
	u.Email = "chuanzegao@163.com"
	u.Avatar = "default.jpg"
	u.Role = 1
	u.Status = 10
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	u.Realname = "管理员"
	o = orm.NewOrm()
	o.Insert(u)
	fmt.Println("insert user end")
}
