package components

import (
	"fmt"
	"library/common"
	"os"
	"strings"
)

type BaseGit struct {
	baseComponents BaseComponents
}

func (c *BaseGit) SetBaseComponents(b BaseComponents) {
	c.baseComponents = b
}
func (c *BaseGit) UpdateRepo(branch string, gitDir string) error {
	if gitDir == "" {
		gitDir = c.baseComponents.GetDeployFromDir()
	}
	if branch == "" {
		branch = "master"
	}
	dotGit := strings.TrimRight(gitDir, "/") + "/.git"
	if _, err := os.Stat(dotGit); err != nil {
		if os.IsNotExist(err) {
			cmds := []string{}
			cmds = append(cmds, fmt.Sprintf("mkdir -p %s ", gitDir))
			cmds = append(cmds, fmt.Sprintf("cd %s ", gitDir))
			cmds = append(cmds, fmt.Sprintf("/usr/bin/env git clone -q %s .", c.baseComponents.project.RepoUrl))
			cmds = append(cmds, fmt.Sprintf("/usr/bin/env git checkout -q %s", branch))
			cmd := strings.Join(cmds, " && ")
			_, err := c.baseComponents.runLocalCommand(cmd)
			return err
		}
	}
	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("cd %s ", gitDir))
	cmds = append(cmds, fmt.Sprintf("/usr/bin/env git fetch --all"))
	cmds = append(cmds, fmt.Sprintf("/usr/bin/env git reset --hard origin/master "))
	cmds = append(cmds, fmt.Sprintf("/usr/bin/env git checkout -q %s ", branch))
	cmds = append(cmds, fmt.Sprintf("/usr/bin/env git fetch -q --all"))
	cmds = append(cmds, fmt.Sprintf("/usr/bin/env git reset -q --hard origin/%s ", branch))
	cmd := strings.Join(cmds, " && ")
	_, err := c.baseComponents.runLocalCommand(cmd)
	return err

}

/**
 * 更新到指定commit版本
 */
func (c *BaseGit) UpdateToVersion() error {
	destination := c.baseComponents.getDeployWorkspace(c.baseComponents.task.LinkId)
	c.UpdateRepo(c.baseComponents.task.Branch, destination)
	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("cd %s ", destination))
	cmds = append(cmds, fmt.Sprintf("/usr/bin/env git reset -q --hard %s ", c.baseComponents.task.CommitId))
	cmd := strings.Join(cmds, " && ")
	_, err := c.baseComponents.runLocalCommand(cmd)
	return err
}

/**
 * 获取分支列表
 */
func (c *BaseGit) GetBranchList() ([]map[string]string, error) {
	history := []map[string]string{}
	destination := c.baseComponents.GetDeployFromDir()
	c.UpdateRepo("master", destination)
	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("cd %s ", destination))
	cmds = append(cmds, "/usr/bin/env git pull -a ")
	cmds = append(cmds, "/usr/bin/env git branch -a ")
	cmd := strings.Join(cmds, " && ")
	s, err := c.baseComponents.runLocalCommand(cmd)
	if err != nil {
		return history, err
	}
	items := strings.Split(s.Result, "\n")
	for _, item := range items {
		item = strings.Trim(item, " ")
		remotePrefix := "remotes/origin/"
		remoteHeadPrefix := "remotes/origin/HEAD"
		if strings.Compare(common.SubString(item, 0, len(remotePrefix)), remotePrefix) == 0 && strings.Compare(common.SubString(item, 0, len(remoteHeadPrefix)), remoteHeadPrefix) != 0 {
			item = common.SubString(item, len(remotePrefix), len(item))
			history = append(history, map[string]string{"id": item, "message": item})
		}
	}
	return history, nil
}

/**
 * 获取提交历史
 *
 */
func (c *BaseGit) GetCommitList(branch string, count int) ([]map[string]string, error) {
	if count == 0 {
		count = 20

	}
	if branch == "" {
		branch = "master"
	}
	history := []map[string]string{}
	destination := c.baseComponents.GetDeployFromDir()
	c.UpdateRepo(branch, destination)
	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("cd %s ", destination))
	cmds = append(cmds, `/usr/bin/env git log -`+common.GetString(count)+` --pretty="%h - %an %s" `)
	cmd := strings.Join(cmds, " && ")
	s, err := c.baseComponents.runLocalCommand(cmd)
	if err != nil {
		return history, err
	}
	items := strings.Split(s.Result, "\n")
	for _, item := range items {
		if strings.Index(item, "-") > -1 {
			commitId := common.SubString(item, 0, strings.Index(item, "-")-1)
			history = append(history, map[string]string{"id": commitId, "message": item})
		}
	}
	return history, nil
}

/**
 * 获取tag记录
 *
 */
func (c *BaseGit) GetTagList(count int) ([]map[string]string, error) {
	if count == 0 {
		count = 20
	}
	c.UpdateRepo("", "")
	history := []map[string]string{}
	destination := c.baseComponents.GetDeployFromDir()
	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("cd %s ", destination))
	cmds = append(cmds, `/usr/bin/env git tag -l `)
	cmd := strings.Join(cmds, " && ")
	s, err := c.baseComponents.runLocalCommand(cmd)
	if err != nil {
		return history, err
	}
	items := strings.Split(s.Result, "\n")
	for _, item := range items {
		history = append(history, map[string]string{"id": item, "message": item})
	}
	return history, nil
}

func (c *BaseGit) DiffBetweenCommits(branch string, commitIdNew string, commitIdOld string) ([]string, error) {
	c.UpdateRepo(branch, "")
	destination := c.baseComponents.GetDeployFromDir()
	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("cd %s ", destination))
	cmds = append(cmds, `/usr/bin/env git diff --name-only  `+commitIdNew+` `+commitIdOld+` `)
	cmd := strings.Join(cmds, " && ")
	s, err := c.baseComponents.runLocalCommand(cmd)
	var files []string
	if err != nil {
		return nil, err
	} else {
		items := strings.Split(s.Result, "\n")
		for _, item := range items {
			if len(item) > 0 {
				files = append(files, item)
			}
		}
		return files, nil
	}
}

func (c *BaseGit) GetLastModifyInfo(branch string, filepath string) (map[string]string, error) {
	destination := c.baseComponents.GetDeployFromDir()
	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("cd %s ", destination))
	cmds = append(cmds, `/usr/bin/env git log -- `+branch+` `+filepath+` | head -3 | tail -2`)
	cmd := strings.Join(cmds, " && ")
	s, err := c.baseComponents.runLocalCommand(cmd)
	if err != nil {
		return nil, err
	} else {
		lines := strings.Split(s.Result, "\n")

		name := common.SubString(lines[0], 8, 100)
		time := common.SubString(lines[1], 8, 100)

		var fileinfo map[string]string
		fileinfo = make(map[string]string)
		fileinfo["name"] = name
		fileinfo["time"] = time
		return fileinfo, nil
	}
}
