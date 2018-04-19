package components

import (
	"fmt"
	"library/common"
	"strings"
)

func (c *BaseComponents) PreDeploy(version string) error {
	tasks := strings.Split(c.project.PreDeploy, "\n")
	if len(tasks) == 0 {
		return nil
	}
	for i, _ := range tasks {
		tasks[i] = strings.Replace(tasks[i], "\n", "", -1)
		tasks[i] = strings.Replace(tasks[i], "\r", "", -1)
	}
	cmds := []string{}
	workspace := strings.TrimRight(c.getDeployWorkspace(version), "/")
	ipsString := strings.Join(c.GetHostIps(), ",")
	ipAndPortString := strings.Join(c.GetAllHost(), ",")
	replaceMap := map[string]string{}
	replaceMap["{WORKSPACE}"] = workspace
	replaceMap["{HOSTS}"] = ipsString
	replaceMap["{HOSTPORT}"] = ipAndPortString
	replaceMap["{ENV}"] = common.GetString(c.project.Level)
	cmds = append(cmds, fmt.Sprintf("cd %s", workspace))
	for _, task := range tasks {
		if strings.Trim(task, " ") != "" {
			taskStr := task
			for k, v := range replaceMap {
				taskStr = strings.Replace(taskStr, k, v, -1)
			}
			if strings.Trim(taskStr, " ") != "" {
				cmds = append(cmds, taskStr)
			}
		}
	}
	cmd := strings.Join(cmds, " && ")
	_, err := c.runLocalCommand(cmd)
	return err

}
func (c *BaseComponents) PostDeploy(version string) error {
	tasks := strings.Split(c.project.PostDeploy, "\n")
	if len(tasks) == 0 {
		return nil
	}
	for i, _ := range tasks {
		tasks[i] = strings.Replace(tasks[i], "\n", "", -1)
		tasks[i] = strings.Replace(tasks[i], "\r", "", -1)
	}
	cmds := []string{}
	workspace := strings.TrimRight(c.getDeployWorkspace(version), "/")
	ipsString := strings.Join(c.GetHostIps(), ",")
	ipAndPortString := strings.Join(c.GetAllHost(), ",")
	replaceMap := map[string]string{}
	replaceMap["{WORKSPACE}"] = workspace
	replaceMap["{HOSTS}"] = ipsString
	replaceMap["{HOSTPORT}"] = ipAndPortString
	replaceMap["{ENV}"] = common.GetString(c.project.Level)
	cmds = append(cmds, fmt.Sprintf("cd %s", workspace))
	for _, task := range tasks {
		if strings.Trim(task, " ") != "" {
			taskStr := task
			for k, v := range replaceMap {
				taskStr = strings.Replace(taskStr, k, v, -1)
			}
			if strings.Trim(taskStr, " ") != "" {
				cmds = append(cmds, taskStr)
			}
		}
	}
	cmd := strings.Join(cmds, " && ")
	_, err := c.runLocalCommand(cmd)
	return err
}

func (c *BaseComponents) getRemotePreReleaseCommand(version string) string {
	tasks := strings.Split(c.project.PreRelease, "\n")
	if len(tasks) == 0 {
		return ""
	}
	for i, _ := range tasks {
		tasks[i] = strings.Replace(tasks[i], "\n", "", -1)
		tasks[i] = strings.Replace(tasks[i], "\r", "", -1)
	}
	cmds := []string{}
	workspace := c.getTargetWorkspace()
	versionDir := c.getReleaseVersionDir(version)
	ipsString := strings.Join(c.GetHostIps(), ",")
	ipAndPortString := strings.Join(c.GetAllHost(), ",")
	replaceMap := map[string]string{}
	replaceMap["{WORKSPACE}"] = workspace
	replaceMap["{VERSION}"] = versionDir
	replaceMap["{HOSTS}"] = ipsString
	replaceMap["{HOSTPORT}"] = ipAndPortString
	replaceMap["{ENV}"] = common.GetString(c.project.Level)
	cmds = append(cmds, fmt.Sprintf("cd %s", versionDir))
	for _, task := range tasks {
		if strings.Trim(task, " ") != "" {
			taskStr := task
			for k, v := range replaceMap {
				taskStr = strings.Replace(taskStr, k, v, -1)
			}
			if strings.Trim(taskStr, " ") != "" {
				cmds = append(cmds, taskStr)
			}
		}
	}
	cmd := strings.Join(cmds, " && ")
	return cmd
}

func (c *BaseComponents) getRemotePostReleaseCommand(version string) string {
	tasks := strings.Split(c.project.PostRelease, "\n")
	if len(tasks) == 0 {
		return ""
	}
	for i, _ := range tasks {
		tasks[i] = strings.Replace(tasks[i], "\n", "", -1)
		tasks[i] = strings.Replace(tasks[i], "\r", "", -1)
	}
	cmds := []string{}
	workspace := c.getTargetWorkspace()
	versionDir := c.getReleaseVersionDir(version)
	ipsString := strings.Join(c.GetHostIps(), ",")
	ipAndPortString := strings.Join(c.GetAllHost(), ",")
	replaceMap := map[string]string{}
	replaceMap["{WORKSPACE}"] = workspace
	replaceMap["{VERSION}"] = versionDir
	replaceMap["{HOSTS}"] = ipsString
	replaceMap["{HOSTPORT}"] = ipAndPortString
	replaceMap["{ENV}"] = common.GetString(c.project.Level)
	cmds = append(cmds, fmt.Sprintf("cd %s", versionDir))
	for _, task := range tasks {
		if strings.Trim(task, " ") != "" {
			taskStr := task
			for k, v := range replaceMap {
				taskStr = strings.Replace(taskStr, k, v, -1)
			}
			if strings.Trim(taskStr, " ") != "" {
				cmds = append(cmds, taskStr)
			}
		}
	}
	cmd := strings.Join(cmds, " && ")
	return cmd
}
func (c *BaseComponents) UpdateRemoteServers(version string) error {
	cmds := []string{}
	// pre-release task
	if c.getRemotePreReleaseCommand(version) != "" {
		cmds = append(cmds, c.getRemotePreReleaseCommand(version))
	}
	if c.GetLinkCommand(version) != "" {
		cmds = append(cmds, c.GetLinkCommand(version))
	}
	if c.getRemotePostReleaseCommand(version) != "" {
		cmds = append(cmds, c.getRemotePostReleaseCommand(version))
	}
	cmd := strings.Join(cmds, " && ")
	_, err := c.runRemoteCommand(cmd, []string{})
	return err
}

func (c *BaseComponents) LastDeploy(version string) error {
	tasks := strings.Split(c.project.LastDeploy, "\n")
	if len(tasks) == 0 {
		return nil
	}
	for i, _ := range tasks {
		tasks[i] = strings.Replace(tasks[i], "\n", "", -1)
		tasks[i] = strings.Replace(tasks[i], "\r", "", -1)
	}
	cmds := []string{}
	replaceMap := map[string]string{}
	replaceMap["{WORKSPACE}"] = strings.TrimRight(c.getDeployWorkspace(version), "/")
	replaceMap["{VERSION}"] = c.getReleaseVersionDir(version)
	ipsString := strings.Join(c.GetHostIps(), ",")
	ipAndPortString := strings.Join(c.GetAllHost(), ",")
	replaceMap["{HOSTS}"] = ipsString
	replaceMap["{HOSTPORT}"] = ipAndPortString
	replaceMap["{PROJECT_ID}"] = common.ToString(c.project.Id)
	replaceMap["{PROJECT_NAME}"] = c.project.Name
	replaceMap["{ENV}"] = common.GetString(c.project.Level)
	if c.task != nil {
		if c.task.Id != 0 {
			replaceMap["{TASK_ID}"] = common.ToString(c.task.Id)
		}
		if c.task.LinkId != "" {
			replaceMap["{TASK_LINKID}"] = common.ToString(c.task.LinkId)
		}
	}
	workspace := strings.TrimRight(c.getDeployWorkspace(version), "/")
	cmds = append(cmds, fmt.Sprintf("cd %s", workspace))
	for _, task := range tasks {
		if strings.Trim(task, " ") != "" {
			taskStr := task
			for k, v := range replaceMap {
				taskStr = strings.Replace(taskStr, k, v, -1)
			}
			if strings.Trim(taskStr, " ") != "" {
				cmds = append(cmds, taskStr)
			}
		}
	}
	cmd := strings.Join(cmds, " && ")
	_, err := c.runLocalCommand(cmd)
	return err

}
