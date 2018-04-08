package api_example

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"log"
	"time"
)

type WalleRes struct {
	Msg  string      `json:"msg,omitempty"`
	Code int         `json:"code,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type Project struct {
	Id             int
	UserId         uint      //必填管理员为0
	Name           string    //必填 项目名字
	Level          int16     // 必填 项目环境1：测试，2：仿真，3：线上
	Status         int16     // 不填
	Version        string    // 不填
	RepoUrl        string    //git地址
	RepoUsername   string    // 不填
	RepoPassword   string    // 不填
	RepoMode       string    // 不填
	RepoType       string    //必填 git/file
	DeployFrom     string    //必填 宿主机存放clone出来的文件
	Excludes       string    //选填 要排除的文件
	ReleaseUser    string    //必填 目标机器用户
	ReleaseTo      string    //必填 目标机器的目录，相当于nginx的root，可直接web访问
	ReleaseLibrary string    //必填 目标机器版本发布库
	Hosts          string    //必填 目标机器列表
	PreDeploy      string    //选填 部署前置任务
	PostDeploy     string    //选填 同步之前任务
	PreRelease     string    //选填 同步之前目标机器执行的任务
	PostRelease    string    //选填 同步之后目标机器执行的任务
	Audit          int16     // 不填
	Ansible        int16     // 不填
	KeepVersionNum int       //选填 线上版本保留数 默认20
	CreatedAt      time.Time // 不填
	UpdatedAt      time.Time // 不填
	P2p            int16     //选填 部署前置任务
	Orgalorg       int16     // 不填
	HostGroup      string    // 不填
	Gzip           int16     // 不填
}

type Task struct {
	ProjectId  int
	Title      string
	CommitId   string
	Branch     string
	PmsBatchId int //选填
	PmsUworkId int //选填
}

var UserAuthKey = "cJIrTa_b2Hnjn6BZkrL8PJkYto2Ael3O"

var WalleUrl = "http://192.168.139.126:8192"

func main() {
	ss, err := getProject("1")
	log.Println(ss)
	log.Println(err)
	pro, err := saveProject(Project{Name: "test"})
	log.Println(pro)
	log.Println(err)
	task, err := getLastTask("1")
	log.Println(task)
	log.Println(err)
	id, err := createTask(1, "test", "test1", "test2", 0, 0)
	log.Println(id)
	log.Println(err)
	res, err := startTask(id)
	log.Println(res)
	log.Println(err)
}

//获取项目
func getProject(projectId string) (map[string]interface{}, error) {
	req := httplib.Get(WalleUrl + "/api/get/conf/get?projectId=" + projectId)
	req.Header("Authorization", "token "+UserAuthKey)
	reqstr, err := req.String()
	if err != nil {
		log.Println("err:", err)
		return map[string]interface{}{}, err
	}

	var walleRes WalleRes
	if err := json.Unmarshal([]byte(reqstr), &walleRes); err != nil {
		return map[string]interface{}{}, errors.New(" walle task fail id=" + projectId + ",msg=" + walleRes.Msg)
	}
	ss := walleRes.Data.(map[string]interface{})
	return ss, nil

}

func saveProject(project Project) (map[string]interface{}, error) {
	req := httplib.Post(WalleUrl + "/api/post/conf/save")
	req.Header("Authorization", "token "+UserAuthKey)
	req.JSONBody(project)
	reqstr, err := req.String()
	if err != nil {
		log.Println("err:", err)
		return map[string]interface{}{}, err
	}
	var walleRes WalleRes
	if err := json.Unmarshal([]byte(reqstr), &walleRes); err != nil {
		return map[string]interface{}{}, errors.New(" msg=" + walleRes.Msg)
	}
	ss := walleRes.Data.(map[string]interface{})
	return ss, nil
}

func getLastTask(projectId string) (map[string]interface{}, error) {
	req := httplib.Get(WalleUrl + "/api/get/task/last?projectId=" + projectId)
	req.Header("Authorization", "token "+UserAuthKey)
	reqstr, err := req.String()
	if err != nil {
		log.Println("err:", err)
		return map[string]interface{}{}, err
	}
	var walleRes WalleRes
	if err := json.Unmarshal([]byte(reqstr), &walleRes); err != nil {
		return map[string]interface{}{}, errors.New(" walle task fail id=" + projectId + ",msg=" + walleRes.Msg)
	}
	ss := walleRes.Data.(map[string]interface{})
	return ss, nil
}

//commit_id 前7位 file方式这里下载路径
func createTask(pro_id int, title string, branch string, commit_id string, pms_batch_id int, pms_uwork_id int) (string, error) {
	req := httplib.Post(WalleUrl + "/api/post/task/save")
	task := Task{
		Title:      title,
		Branch:     branch,
		CommitId:   commit_id,
		ProjectId:  pro_id,
		PmsBatchId: pms_batch_id,
		PmsUworkId: pms_uwork_id,
	}
	req.Header("Authorization", "token "+UserAuthKey)
	req.JSONBody(task)
	reqstr, err := req.String()
	if err != nil {
		log.Println("err:", err)
		return "", errors.New("createTask api fail")
	}
	log.Println(reqstr)
	var walleRes WalleRes
	if err := json.Unmarshal([]byte(reqstr), &walleRes); err != nil {
		log.Println("err:", err)
		return "", errors.New("createTask  json fail")
	}
	ss := walleRes.Data.(map[string]interface{})
	return fmt.Sprintf("%v", ss["Id"]), nil

}
func startTask(taskId string) (string, error) {
	req := httplib.Get(WalleUrl + "/api/get/walle/release?taskId=" + taskId)
	req.Header("Authorization", "token "+UserAuthKey)
	reqstr, err := req.String()
	if err != nil {
		log.Println("err:", err)
		return "timeout", nil
	}
	var walleRes WalleRes
	if err := json.Unmarshal([]byte(reqstr), &walleRes); err != nil {
		log.Println("err:", err)
		return "", errors.New("startTask  json fail")
	}
	if walleRes.Code != 0 {
		return "fail", errors.New(" walle task fail id=" + taskId + ",msg=" + walleRes.Msg)
	} else {
		return "success", nil
	}
}
