package components

import ()
import (
	"fmt"
	"github.com/astaxie/beego"
	"library/common"
	"library/p2p/init_sever"
	"models"
	"strings"
	"time"
)

/**
 * 初始化宿主机部署工作空间
 *
 */
func (c *BaseComponents) InitLocalWorkspace(version string) error {
	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("cp -rf %s %s ", c.GetDeployFromDir(), c.getDeployWorkspace(version)))
	cmd := strings.Join(cmds, " && ")
	_, err := c.runLocalCommand(cmd)
	return err
}

/**
 * 目标机器的版本库初始化
 */
func (c *BaseComponents) InitRemoteVersion(version string) error {
	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("mkdir -p %s ", c.getReleaseVersionDir(version)))
	cmd := strings.Join(cmds, " && ")
	_, err := c.runRemoteCommand(cmd, []string{})
	return err
}

/**
 * 将多个文件/目录传输到指定的多个目标机

 */
func (c *BaseComponents) CopyFiles() error {
	err := c.packageFiles()
	if err != nil {
		return err
	}
	version := c.task.LinkId
	src := c.getDeployPackagePath(version)
	dest := c.getReleaseVersionPackage(version)
	if c.project.P2p == 1 {
		_, err := c.copyFilesByP2p(version, src, dest, []string{})
		if err != nil {
			return err
		}
	} else {
		_, err := c.copyFilesBySftp(src, dest, []string{})
		if err != nil {
			beego.Info(err)
			return err
		}
	}
	err = c.unpackageFiles()
	if err != nil {
		return err
	}
	return nil
}

/**
 * 打软链

 */
func (c *BaseComponents) GetLinkCommand(version string) string {
	user := c.project.ReleaseUser
	project := c.GetGitProjectName(c.project.RepoUrl)
	linkFrom := c.getReleaseVersionDir(version)
	currentTmp := fmt.Sprintf("%s/%s/current-%s.tmp", strings.TrimRight(c.project.ReleaseLibrary, "/"), project, project)

	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("cd %s ", linkFrom))

	if c.project.ReleaseType == models.RELEASE_TYPE_SOFTLINK {
		cmds = append(cmds, fmt.Sprintf("ln -sfn %s %s ", linkFrom, currentTmp))
		cmds = append(cmds, fmt.Sprintf("chown -h %s %s ", user, currentTmp))
		cmds = append(cmds, fmt.Sprintf("mv -fT %s %s ", currentTmp, c.project.ReleaseTo))
	} else {
		cmds = append(cmds, fmt.Sprintf("cp -r %s %s ", linkFrom, currentTmp))
		cmds = append(cmds, fmt.Sprintf("chown -h %s %s ", user, currentTmp))
		trashSuffix := "_trash_" + time.Now().Format("20060102_150402")
		cmds = append(cmds, fmt.Sprintf("mv %s %s;mv -fT %s %s ", c.project.ReleaseTo, c.project.ReleaseTo+trashSuffix, currentTmp, c.project.ReleaseTo))
	}
	cmd := strings.Join(cmds, " && ")
	return cmd
}

/**
 * 解包文件
 */
func (c *BaseComponents) unpackageFiles() error {
	version := c.task.LinkId
	releasePath := c.getReleaseVersionDir(version)
	releasePackage := c.getReleaseVersionPackage(version)
	unTarparameter := "-xzf"
	if c.project.Status == 0 {
		unTarparameter = "-xf"
	}
	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("cd %s ", releasePath))
	//兼容docker
	if beego.BConfig.RunMode == "docker" {
		cmds = append(cmds, fmt.Sprintf("tar %s %s -C %s", unTarparameter, releasePackage, releasePath))
	} else {
		cmds = append(cmds, fmt.Sprintf("tar --preserve-permissions --touch --no-same-owner %s %s -C %s", unTarparameter, releasePackage, releasePath))
	}
	cmd := strings.Join(cmds, " && ")
	_, err := c.runRemoteCommand(cmd, []string{})
	return err
}

/**
 * 打包文件
 */
func (c *BaseComponents) packageFiles() error {
	version := c.task.LinkId
	packagePath := c.getDeployPackagePath(version)
	tarparameter := "-czf"
	if c.project.Status == 0 {
		tarparameter = "-cf"
	}
	cmds := []string{}
	beego.Info(strings.TrimRight(c.getDeployWorkspace(version), "/"))
	cmds = append(cmds, fmt.Sprintf("cd %s", strings.TrimRight(c.getDeployWorkspace(version), "/")))
	commandFiles := "."
	if beego.BConfig.RunMode == "docker" {
		cmds = append(cmds, fmt.Sprintf("tar %s  %s %s %s", c.excludes(version), tarparameter, packagePath, commandFiles))

	} else {
		cmds = append(cmds, fmt.Sprintf("tar %s --preserve-permissions %s %s %s", c.excludes(version), tarparameter, packagePath, commandFiles))
	}
	cmd := strings.Join(cmds, " && ")
	_, err := c.runLocalCommand(cmd)
	return err
}

/**
 * 打包文件
 */
func (c *BaseComponents) excludes(version string) string {
	excludesArr := []string{}
	excludes := strings.Split(c.project.Excludes, "\n")
	for _, exclude := range excludes {
		exclude = strings.Trim(exclude, " ")
		exclude = strings.Trim(exclude, "\r")
		if exclude != "" {
			if !common.InList(exclude, excludesArr) {
				excludesArr = append(excludesArr, exclude)
			}

		}
	}
	cmds := []string{}
	for _, e := range excludesArr {
		cmds = append(cmds, fmt.Sprintf("--exclude=%s", "'"+e+"'"))
	}
	cmd := strings.Join(cmds, " ")
	return cmd

}

/**
 * 清理空间
 */
func (c *BaseComponents) CleanUpLocal(version string) error {
	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("rm -rf %s ", c.getDeployWorkspace(version)))
	cmds = append(cmds, fmt.Sprintf("rm -f %s ", c.getDeployPackagePath(version)))
	cmd := strings.Join(cmds, "&&")
	_, err := c.runLocalCommand(cmd)
	return err
}

/**
 * 清理远端空间
 */
func (c *BaseComponents) CleanUpReleasesVersion(versions []string) error {
	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("cd %s ", c.getReleaseVersionDir("")))
	cmds = append(cmds, fmt.Sprintf("rm -f %s/*.tar.gz", strings.TrimRight(c.getReleaseVersionDir(""), "/")))
	for _, version := range versions {
		if version != "" {
			cmds = append(cmds, fmt.Sprintf("rm -rf %s/%s", strings.TrimRight(c.getReleaseVersionDir(""), "/"), version))
		}
	}
	cmd := strings.Join(cmds, " && ")
	_, err := c.runRemoteCommand(cmd, []string{})
	return err
}

/**
 * 测试ssh连接
 *
 */
func (c *BaseComponents) TestSsh() error {
	cmd := "id"
	_, err := c.runRemoteCommand(cmd, []string{})
	return err
}

/**
 * 检测php用户是否具有目标机release目录读写权限
 *
 */
func (c *BaseComponents) TestReleaseDir() error {
	temDir := "detection" + time.Now().Format("2006-01-02 15:04:05")
	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("mkdir -p %s", c.getReleaseVersionDir(temDir)))
	cmds = append(cmds, fmt.Sprintf("rm -rf  %s", c.getReleaseVersionDir(temDir)))
	cmd := strings.Join(cmds, "&&")
	_, err := c.runRemoteCommand(cmd, []string{})
	return err
}

/**
 * 发送p2p客户端并启动
 *
 */
func (c *BaseComponents) SendP2pAgent(dirAgentPath string, destPath string) error {
	agentFile := "agent_main"
	cmds2 := []string{}
	cmds2 = append(cmds2, fmt.Sprintf("ps -ef |grep %s| grep -v grep  |awk '{print $2}' |xargs kill -9", agentFile))
	cmd2 := strings.Join(cmds2, " && ")
	c.runRemoteCommand(cmd2, []string{})

	cmds1 := []string{}
	cmds1 = append(cmds1, fmt.Sprintf("mkdir -p  %s", strings.TrimRight(destPath, "/")+"/src"))
	cmds1 = append(cmds1, fmt.Sprintf("rm -rf   %s/src/%s", strings.TrimRight(destPath, "/"), agentFile))
	cmd1 := strings.Join(cmds1, "&&")
	_, err := c.runRemoteCommand(cmd1, []string{})
	if err != nil {
		return err
	}
	_, err = c.copyFilesBySftp(strings.TrimRight(dirAgentPath, "/")+"/"+agentFile, strings.TrimRight(destPath, "/")+"/src/"+agentFile, []string{})
	if err != nil {
		return err
	}
	agentFileConf := "agent.json"
	_, err = c.copyFilesBySftp(strings.TrimRight(dirAgentPath, "/")+"/"+agentFileConf, strings.TrimRight(destPath, "/")+"/src/"+agentFileConf, []string{})
	if err != nil {
		return err
	}
	controlFile := "control"
	_, err = c.copyFilesBySftp(strings.TrimRight(dirAgentPath, "/")+"/"+controlFile, strings.TrimRight(destPath, "/")+"/"+controlFile, []string{})
	if err != nil {
		return err
	}
	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("chmod 777 -R %s/* ", strings.TrimRight(destPath, "/")))
	cmds = append(cmds, fmt.Sprintf("cd %s/", strings.TrimRight(destPath, "/")))
	cmds = append(cmds, "./control start")
	cmd := strings.Join(cmds, "&&")
	_, err = c.runRemoteCommand(cmd, []string{})
	if err != nil {
		return err
	}
	init_sever.P2pSvc.CheckAllClientIp([]string{})
	return nil
}

func (c *BaseComponents) GetExecFlush() error {
	projectFile := strings.TrimRight(c.project.ReleaseTo, "/") + "/" + "Conf/extend.php"
	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("cd %s", "/data/www/static/prod/static"))
	cmds = append(cmds, "git log | head -n 1 | awk '{print  substr($2,0,16)}'")
	cmd := strings.Join(cmds, "&&")
	s, err := c.runLocalCommand(cmd)
	if err != nil {
		return err
	}
	sha := s.Result + "_" + common.GetString(time.Now().Unix())
	cmds1 := []string{}
	cmds1 = append(cmds1, fmt.Sprintf("sed -i -r 's/ts=[0-9a-f_]*/ts=%s/g' %s", sha, projectFile))
	cmd1 := strings.Join(cmds1, "&&")
	_, err = c.runRemoteCommand(cmd1, []string{})
	return err
}
func (c *BaseComponents) GetGitLog() error {
	cmds1 := []string{}
	cmds1 = append(cmds1, fmt.Sprintf("cd %s && git branch  | grep  \"*\" && git log -1", c.project.ReleaseTo))
	cmd1 := strings.Join(cmds1, "&&")
	_, err := c.runRemoteCommand(cmd1, []string{})
	return err
}
func (c *BaseComponents) GetGitPull() error {
	cmds1 := []string{}
	cmds1 = append(cmds1, fmt.Sprintf("cd %s && git pull && rm Runtime/* -rf", c.project.ReleaseTo))
	cmd1 := strings.Join(cmds1, "&&")
	_, err := c.runRemoteCommand(cmd1, []string{})
	return err
}
func (c *BaseComponents) StartP2pAgent(ips []string, destPath string) error {
	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("cd %s/", strings.TrimRight(destPath, "/")))
	cmds = append(cmds, "./control start")
	cmd := strings.Join(cmds, "&&")
	_, err := c.runRemoteCommand(cmd, ips)
	return err
}
