package wallecontrollers

import (
	"controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/orm"
	"library/common"
	"library/components"
	"models"
	"strings"
	"time"
)

type ReleaseController struct {
	controllers.BaseController
}

func (c *ReleaseController) Get() {
	if c.Task == nil || c.Task.Id == 0 {
		c.SetJson(1, nil, "Parameter error")
		return
	}
	c.Project, _ = models.GetProjectById(c.Task.ProjectId)
	if c.Project == nil || c.Project.Id == 0 {
		c.SetJson(1, nil, "Parameter error")
		return
	}
	if c.User == nil || c.User.Id == 0 {
		c.SetJson(2, nil, "not login")
		return
	}
	//不是自己的项目不允许上线(除非是管理员项目)
	if c.User.Id != int(c.Task.UserId) && c.User.Id != 1 && c.Task.UserId != 1 {
		c.SetJson(1, nil, "not such uer")
		return
	}
	//上线成功 以及审核失败不允许上线
	if c.Task.Status == 2 || c.Task.Status == 3 {
		c.SetJson(1, nil, "此项目已完成")
		return
	}
	//正在上线的不允许上线
	if c.Task.IsRun == 1 {
		c.SetJson(1, nil, "此项目正在上线中")
		return
	}
	//删除上线日志记录
	o := orm.NewOrm()
	o.Raw("DELETE FROM `record` WHERE `task_id`= ? ", c.Task.Id).Exec()
	if c.Task.Action == 0 {
		//生成版本号
		c.makeVersion()
		go func() {
			err := c.releaseHandling()
			if err != nil {
				models.AddTaskErrLog(&models.TaskErrLog{
					ErrInfo: err.Error(),
					TaskId:  c.Task.Id,
				})
			}
		}()
	} else {
		go c.rollBackHandling()
	}
	c.SetJson(0, nil, "")
	return

}
func (c *ReleaseController) makeVersion() {
	c.Task.IsRun = 1
	version := time.Now().Format("20060102-150405")
	c.Task.LinkId = version
	c.Task.UpdatedAt = time.Now()
	models.UpdateTaskById(c.Task)
}

//回滚任务
func (c *ReleaseController) rollBackHandling() error {
	s := components.BaseComponents{}
	s.SetProject(c.Project)
	s.SetTask(c.Task)
	g := components.BaseGit{}
	g.SetBaseComponents(s)
	err := s.UpdateRemoteServers(c.Task.LinkId)
	if err != nil {
		c.failHandling(&s)
		return err
	}
	err = c.changeReleaseData()
	if err != nil {
		return err
	}
	return nil
}

//普通上线任务
func (c *ReleaseController) updateRecord(action int) error {
	o := orm.NewOrm()
	_, err := o.Raw("UPDATE `record` SET `action`= ?  WHERE`task_id` = ? and action=0", action, c.Task.Id).Exec()
	if err != nil {
		return err
	}
	return nil
}

//普通上线任务
func (c *ReleaseController) releaseHandling() error {
	s := components.BaseComponents{}
	s.SetProject(c.Project)
	s.SetTask(c.Task)

	err := s.InitLocalWorkspace(c.Task.LinkId)
	c.updateRecord(10)
	if err != nil {
		c.failHandling(&s)
		return err
	}
	err = s.InitRemoteVersion(c.Task.LinkId)
	c.updateRecord(10)
	if err != nil {
		c.failHandling(&s)
		return err
	}
	err = s.PreDeploy(c.Task.LinkId)
	c.updateRecord(20)
	if err != nil {
		c.failHandling(&s)
		return err
	}
	if c.Project.RepoType == "git" {
		g := components.BaseGit{}
		g.SetBaseComponents(s)
		err = g.UpdateToVersion()
		c.updateRecord(30)
		if err != nil {
			c.failHandling(&s)
			return err
		}
	} else if c.Project.RepoType == "file" {
		f := components.BaseFile{}
		f.SetBaseComponents(s)
		err = f.UpdateToVersion()
		c.updateRecord(30)
		if err != nil {
			c.failHandling(&s)
			return err
		}
	}

	err = s.PostDeploy(c.Task.LinkId)
	c.updateRecord(40)
	if err != nil {
		c.failHandling(&s)
		return err
	}
	err = s.CopyFiles()
	c.updateRecord(50)
	if err != nil {
		c.failHandling(&s)
		return err
	}
	err = s.UpdateRemoteServers(c.Task.LinkId)
	c.updateRecord(60)
	if err != nil {
		c.failHandling(&s)
		return err
	}
	//这里实际发布已完成 (后置本地脚本任务,)
	err = s.LastDeploy(c.Task.LinkId)
	if err != nil {
		c.failHandling(&s)
		return err
	}

	err = s.CleanUpLocal(c.Task.LinkId)
	c.updateRecord(100)
	if err != nil {
		return err
	}

	err = c.changeReleaseData()
	if err != nil {
		return err
	}
	go c.callUwork("1")
	return nil

}
func (c *ReleaseController) callUwork(isFail string) {
	if c.Task.UserId == 1 {
		url := beego.AppConfig.String("uworkHost") + "uwork/admin/walleCallback?task_id=" + common.GetString(c.Task.Id) + "&res=" + isFail
		if beego.BConfig.RunMode != "prod" {
			url = `http://192.168.149.61:8092/`
		}
		req1 := httplib.Get(url)
		rspUrl2, _ := req1.String()
		beego.Info(rspUrl2)
	}
}
func (c *ReleaseController) changeReleaseData() error {
	//对于回滚的任务不记录线上版本
	if c.Task.Action == 0 {
		c.Task.ExLinkId = c.Project.Version
	}
	//判断是否为第一次任务，或者为回滚任务
	if c.Project.Version == "" || c.Task.Action == 1 {
		c.Task.EnableRollback = 0
	}
	c.Task.Status = 3
	c.Task.IsRun = 0
	c.Task.UpdatedAt = time.Now()
	err := models.UpdateTaskById(c.Task)
	if err != nil {
		return err
	}
	err = c.enableRollBack()
	if err != nil {
		return err
	}
	// 记录当前线上版本（软链）回滚则是回滚的版本，上线为新版本
	c.Project.Version = c.Task.LinkId
	err = models.UpdateProjectById(c.Project)
	if err != nil {
		return err
	}
	return nil
}

func (c *ReleaseController) enableRollBack() error {
	var ids []orm.Params
	o := orm.NewOrm()
	s, err := o.Raw("SELECT id FROM task WHERE `status`=3 and project_id = ? and  `enable_rollback`=1 ORDER BY id DESC LIMIT ?", c.Task.ProjectId, c.Project.KeepVersionNum).Values(&ids)
	if s > 0 && err == nil {
		idStrs := []string{}
		for _, id := range ids {
			idstr := common.GetString(id["id"])
			idStrs = append(idStrs, idstr)
		}
		sqlIn := strings.Join(idStrs, ",")
		var versionsRes []orm.Params
		s1, err := o.Raw("SELECT link_id FROM task WHERE `enable_rollback`=1 and `id` not in ("+sqlIn+") and  project_id = ? ", c.Task.ProjectId).Values(&versionsRes)
		if s1 > 0 && err == nil {
			var versions []string
			for _, version := range versionsRes {
				versions = append(versions, common.GetString(version["link_id"]))
			}
			//这里查找需要设置不可回滚的版本 进行清除操作
			s := components.BaseComponents{}
			s.SetProject(c.Project)
			s.SetTask(c.Task)
			s.CleanUpReleasesVersion(versions)
		}
		_, err = o.Raw("UPDATE `task` SET `enable_rollback`='0' WHERE`id` not in ("+sqlIn+") and  project_id = ? and  `enable_rollback`=1 ", c.Task.ProjectId).Exec()
		if err != nil {
			return err
		}
		return nil
	} else {
		return err
	}
}

//上线失败处理
func (c *ReleaseController) failHandling(co *components.BaseComponents) {
	//修改状态
	c.Task.Status = 4
	c.Task.IsRun = 0
	c.Task.UpdatedAt = time.Now()
	models.UpdateTaskById(c.Task)
	//清理本地版本库
	co.CleanUpLocal(c.Task.LinkId)
	go c.callUwork("0")
}
