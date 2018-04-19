package components

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/gaoyue1989/sshexec"
	"library/common"
	"library/ssh"
	"models"
	"regexp"
	"strings"
	"time"
)

const SSHTIMEOUT = 3600
const SSHWorker = 10

const SSHREMOTETIMEOUT = 600

type BaseComponents struct {
	name    string
	project *models.Project
	task    *models.Task
}

func (c *BaseComponents) SetProject(project *models.Project) {
	c.project = project

}
func (c *BaseComponents) SetTask(task *models.Task) {
	c.task = task
}

/**
* 执行本地宿主机命令
 */
func (c *BaseComponents) runLocalCommand(command string) (sshexec.ExecResult, error) {
	id := c.SaveRecord(command)
	s, err := gopubssh.CommandLocal(command, SSHTIMEOUT)
	ss, _ := json.Marshal(s)
	go c.LogTaskCommond(string(ss))
	//获取执行时间
	duration := common.GetInt(s.EndTime.Sub(s.StartTime).Seconds())
	createdAt := int(s.StartTime.Unix())
	status := 1
	if s.Error != nil {
		status = 0
	}
	c.SaveRecordRes(id, duration, createdAt, status, s)
	return s, err

}

/**
* 执行远端目标机命令
 */
func (c *BaseComponents) runRemoteCommand(command string, hosts []string) ([]sshexec.ExecResult, error) {
	if len(hosts) == 0 {
		hostsInfo:=c.GetHosts()
		for _, info := range hostsInfo {
			hosts=append(hosts,info.AllHost)
		}
	}
	id := c.SaveRecord(command)
	start := time.Now()
	createdAt := int(start.Unix())
	sshExecAgent := sshexec.SSHExecAgent{}
	sshExecAgent.Worker = SSHWorker
	sshExecAgent.TimeOut = time.Duration(SSHREMOTETIMEOUT) * time.Second
	s, err := sshExecAgent.SshHostByKey(hosts, c.project.ReleaseUser, command)
	ss, _ := json.Marshal(s)
	go c.LogTaskCommond(string(ss))
	//获取执行时间
	duration := common.GetInt(time.Now().Sub(start).Seconds())

	status := 1
	if err != nil {
		status = 0
	}
	c.SaveRecordRes(id, duration, createdAt, status, s)
	return s, err

}

/**
* 执行远端传输文件
 */
func (c *BaseComponents) copyFilesBySftp(src string, dest string, hosts []string) ([]sshexec.ExecResult, error) {
	if len(hosts) == 0 {
		hostsInfo:=c.GetHosts()
		for _, info := range hostsInfo {
			hosts=append(hosts,info.AllHost)
		}
	}
	id := c.SaveRecord("Transfer")
	start := time.Now()
	createdAt := int(start.Unix())
	sshExecAgent := sshexec.SSHExecAgent{}
	sshExecAgent.Worker = SSHWorker
	sshExecAgent.TimeOut = time.Duration(SSHREMOTETIMEOUT) * time.Second
	s, err := sshExecAgent.SftpHostByKey(hosts, c.project.ReleaseUser, src, dest)
	ss, _ := json.Marshal(s)
	go c.LogTaskCommond(string(ss))
	//获取执行时间
	duration := common.GetInt(time.Now().Sub(start).Seconds())
	status := 1
	if err != nil {
		status = 0
	}
	c.SaveRecordRes(id, duration, createdAt, status, s)
	return s, err

}

/**
* 执行远端传输文件 p2p方式
 */
func (c *BaseComponents) copyFilesByP2p(id string, src string, dest string, hosts []string) ([]sshexec.ExecResult, error) {
	start := time.Now()
	rid := c.SaveRecord("Transfer by p2p")
	createdAt := int(start.Unix())
	if len(hosts) == 0 {
		hostsInfo:=c.GetHosts()
		for _, info := range hostsInfo {
			hosts=append(hosts,info.Ip)
		}
	}
	s, err := gopubssh.TransferByP2p(id, hosts, c.project.ReleaseUser, src, dest, SSHREMOTETIMEOUT)
	ss, _ := json.Marshal(s)
	go c.LogTaskCommond(string(ss))
	//获取执行时间
	duration := common.GetInt(time.Now().Sub(start).Seconds())

	status := 1
	if err != nil {
		status = 0
	}
	c.SaveRecordRes(rid, duration, createdAt, status, s)
	return s, err

}


type HostInfo struct {
	Ip    string
	Group    int
	Port  int
	AllHost string
}


/**
 * 获取host
 */
func (c *BaseComponents) GetHosts() []HostInfo {
	hostsStr := c.project.Hosts
	if c.task != nil && c.task.Hosts != "" {
		hostsStr = c.task.Hosts
	}
	//获取ip
	reg := regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)\.(\d+)`)
	hosts := reg.FindAll([]byte(hostsStr), -1)
	res := []HostInfo{}
	for _, host := range hosts {
		isInList:=false
		for _, r := range res {
			if r.Ip==string(host){
				isInList=true
			}
		}
		if !isInList {
			res = append(res, HostInfo{Ip:string(host),Port:22})
		}
	}
	//格式化端口号
	reg1 := regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)\.(\d+)\:(\d+)`)
	hosts1 := reg1.FindAll([]byte(hostsStr), -1)
	for _, host := range hosts1 {
		ip:=strings.Split(string(host),":")[0]
		port:=strings.Split(string(host),":")[1]
		for i, r := range res {
			if r.Ip==ip{
				res[i].Port=common.GetInt(port)
			}
		}
	}
	//格式化端口号
	reg2 := regexp.MustCompile(`(\d+)\#(\d+)\.(\d+)\.(\d+)\.(\d+)`)
	hosts2 := reg2.FindAll([]byte(hostsStr), -1)
	for _, host := range hosts2 {
		ip:=strings.Split(string(host),"#")[1]
		group:=strings.Split(string(host),"#")[0]
		for i, r := range res {
			if r.Ip==ip{
				res[i].Group=common.GetInt(group)
			}
		}
	}
	for i,  r:= range res {
		res[i].AllHost=r.Ip+":"+common.GetString(r.Port)
	}
	return res
}
/**
 * 获取host ip
 */
func (c *BaseComponents) GetHostIps() []string {
	hosts:=[]string{}
	hostsInfo:=c.GetHosts()
	for _, info := range hostsInfo {
		hosts=append(hosts,info.Ip)
	}
	return hosts
}
/**
 * 获取host ip加端口
 */
func (c *BaseComponents) GetAllHost() []string {
	hosts:=[]string{}
	hostsInfo:=c.GetHosts()
	for _, info := range hostsInfo {
		hosts=append(hosts,info.AllHost)
	}
	return hosts
}
func (c *BaseComponents) GetGroupHost() map[int]string {
	hosts:=map[int]string{}
	hostsInfo:=c.GetHosts()
	beego.Info(hostsInfo)
	for _, info := range hostsInfo {
		hosts[info.Group]=info.Ip+":"+common.GetString(info.Port)+"\r\n"
	}
	return hosts
}
/**
 * 获取环境
 */
func (c *BaseComponents) getEnv() string {
	if c.project.Level == 1 {
		return "test"
	}
	if c.project.Level == 2 {
		return "simu"
	}
	if c.project.Level == 3 {
		return "prod"
	}
	return "unknow"
}

/**
 * 拼接宿主机的部署隔离工作空间
 * {deploy_from}/{env}/{project}-YYmmdd-HHiiss
 */
func (c *BaseComponents) getDeployWorkspace(version string) string {
	from := c.project.DeployFrom
	env := c.getEnv()
	project := c.GetGitProjectName(c.project.RepoUrl)
	return fmt.Sprintf("%s/%s/%s-%s", strings.TrimRight(from, "/"), strings.TrimRight(env, "/"), project, version)
}

/**
 * 获取传输宿主机tar文件路径
 *
 * {deploy_from}/{env}/{project}-YYmmdd-HHiiss.tar.gz
 */
func (c *BaseComponents) getDeployPackagePath(version string) string {
	return fmt.Sprintf("%s.tar.gz", c.getDeployWorkspace(version))
}

/**
 * 拼接宿主机的仓库目录
 * {deploy_from}/{env}/{project}
 */
func (c *BaseComponents) GetDeployFromDir() string {

	from := c.project.DeployFrom
	env := c.getEnv()
	project := c.GetGitProjectName(c.project.RepoUrl)
	return fmt.Sprintf("%s/%s/%s", strings.TrimRight(from, "/"), strings.TrimRight(env, "/"), project)
}

/**
 * 获取目标机要发布的目录
 * {webroot}
 */
func (c *BaseComponents) getTargetWorkspace() string {
	return strings.TrimRight(c.project.ReleaseTo, "/")
}

/**
 * 拼接目标机要发布的目录
 * {release_library}/{project}/{version}
 */
func (c *BaseComponents) getReleaseVersionDir(version string) string {
	return fmt.Sprintf("%s/%s/%s", strings.TrimRight(c.project.ReleaseLibrary, "/"), c.GetGitProjectName(c.project.RepoUrl), version)
}

/**
 * 拼接目标机要发布的打包文件路径
 * {release_library}/{project}/{version}
 */
func (c *BaseComponents) getReleaseVersionPackage(version string) string {
	return fmt.Sprintf("%s.tar.gz", c.getReleaseVersionDir(version))
}

//根据git地址获取项目名字
func (c *BaseComponents) GetGitProjectName(gitUrl string) string {
	s := strings.Split(gitUrl, "/")
	sname := s[len(s)-1]
	snames := strings.Split(sname, `.git`)
	if snames[0] == "" {
		return "filedir"
	}
	return snames[0]
}

func (c *BaseComponents) LogTaskCommond(value interface{}) {

	////设置日志
	//fn := "logs/task_log/task-" + time.Now().Format("20060102") + ".log"
	//if _, err := os.Stat(fn); err != nil {
	//	if os.IsNotExist(err) {
	//		os.Create(fn)
	//	}
	//}
	//log := logs.NewLogger(1)
	//log.SetLogger("file", `{"filename":"` + fn + `"}`)
	//log.Info("---------------------------------")
	//log.Info("id:%d > %s\n", c.task.Id, value)
	//log.Info("---------------------------------")
}
func (c *BaseComponents) SaveRecord(command string) int {
	re := models.Record{}
	re.Command = command
	if c.task == nil || c.task.Id == 0 {
		re.TaskId = -99
		re.UserId = 0
	} else {
		re.TaskId = int64(c.task.Id)
		re.UserId = c.task.UserId
	}
	re.Status = 1
	id, err := models.AddRecord(&re)
	if err != nil {
		beego.Error(err)
	}
	return int(id)
}
func (c *BaseComponents) SaveRecordRes(id int, duration int, createdAt int, status int, value interface{}) {
	beego.Info(value)
	if duration < 0 {
		duration = 0
	}
	re, err := models.GetRecordById(id)
	if err != nil {
		beego.Error(err)
		return
	}
	re.Duration = duration
	sResult, _ := json.Marshal(value)
	re.CreatedAt = createdAt
	re.Memo = string(sResult)
	re.Status = int16(status)
	err = models.UpdateRecordById(re)
	if err != nil {
		beego.Error(err)
	}
}

/**
 * 清理项目目录
 *
 */
func (c *BaseComponents) RemoveLocalProjectWorkspace() error {
	gitDir := c.GetDeployFromDir()
	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("rm -rf  %s ", gitDir))
	cmd := strings.Join(cmds, "&&")
	_, err := c.runLocalCommand(cmd)
	return err
}
