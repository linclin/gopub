package gopubssh

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gaoyue1989/sshexec"
	"library/p2p/init_sever"
	"library/p2p/server"
	"os/exec"
	"path/filepath"
	"time"
)

func CommandLocal(cmd string, to int) (sshexec.ExecResult, error) {
	timeout := time.After(time.Duration(to) * time.Second)
	execResultCh := make(chan *sshexec.ExecResult, 1)
	go func() {
		execResult := LocalExec(cmd)
		execResultCh <- &execResult
	}()
	select {
	case res := <-execResultCh:
		sres := *res
		errorText := ""
		if sres.Error != nil {
			errorText += " commond  exec error.\n" + "rsult info :" + sres.Result + "\nerror info :" + sres.Error.Error()
		}
		if errorText != "" {
			return sres, errors.New(errorText)
		} else {
			return sres, nil
		}

	case <-timeout:
		return sshexec.ExecResult{Command: cmd, Error: errors.New("cmd time out")}, errors.New("cmd time out")
	}

}
func LocalExec(cmd string) sshexec.ExecResult {
	execResult := sshexec.ExecResult{}
	execResult.StartTime = time.Now()
	execResult.Command = cmd
	execCommand := exec.Command("/bin/bash", "-c", cmd)

	var b bytes.Buffer

	execCommand.Stdout = &b
	var b1 bytes.Buffer
	execCommand.Stderr = &b1
	err := execCommand.Run()
	if err != nil {
		execResult.Error = err
		execResult.Result = b1.String()
		return execResult
	} else {
		execResult.EndTime = time.Now()
		execResult.Result = b.String()
		return execResult
	}
}

func TransferByP2p(id string, hosts []string, user string, localFilePath string, remoteFilePath string, to int) ([]sshexec.ExecResult, error) {
	returnResult := make([]sshexec.ExecResult, len(hosts))
	timeout := time.After(time.Duration(to) * time.Second)
	//创建传输任务
	s := server.CreateTask{ID: id, DispatchFiles: []string{localFilePath}, DestIPs: hosts}
	init_sever.P2pSvc.CreateTaskNoHttp(&s)
	taskInfoCh := make(chan *server.TaskInfo, 1)

	go func() {
		for {
			ss, _ := init_sever.P2pSvc.QueryTaskNoHttp(id)
			if ss.Status == server.TaskCompleted.String() {
				taskInfoCh <- ss
				break
			} else if ss.Status == server.TaskFailed.String() {
				taskInfoCh <- ss
				break
			}
			time.Sleep(100 * time.Millisecond)

		}
	}()
	select {
	case res := <-taskInfoCh:
		if res.Status == server.TaskCompleted.String() {
			e := sshexec.ExecResult{}
			for ip, DispatchInfo := range res.DispatchInfos {
				e.LocalFilePath = localFilePath
				e.RemoteFilePath = remoteFilePath
				e.StartTime = DispatchInfo.StartedAt
				e.EndTime = DispatchInfo.FinishedAt
				e.Host = ip
				returnResult = append(returnResult, e)
			}
			err := TransP2pReName(id, hosts, user, localFilePath, remoteFilePath, 30)
			return returnResult, err
		} else {
			for ip, DispatchInfo := range res.DispatchInfos {
				e := sshexec.ExecResult{}
				if DispatchInfo.Status != server.TaskCompleted.String() {
					e.LocalFilePath = localFilePath
					e.RemoteFilePath = remoteFilePath
					e.StartTime = DispatchInfo.StartedAt
					e.Error = errors.New("p2p transfer error")
					e.Host = ip
				} else {
					e.LocalFilePath = localFilePath
					e.RemoteFilePath = remoteFilePath
					e.StartTime = DispatchInfo.StartedAt
					e.EndTime = DispatchInfo.FinishedAt
					e.Host = ip
				}
				returnResult = append(returnResult, e)
			}
			return returnResult, errors.New("p2p transfer error")

		}

	case <-timeout:
		return returnResult, errors.New("p2p time out")
	}
}

func TransP2pReName(id string, hosts []string, user string, localFilePath string, remoteFilePath string, to int) error {
	fileName := filepath.Base(localFilePath)
	filePath := init_sever.P2pSvc.Cfg.DownDir
	oldFile := filePath + fileName
	cmd := fmt.Sprintf("mv -f %s %s", oldFile, remoteFilePath)
	sshExecAgent := sshexec.SSHExecAgent{}
	sshExecAgent.Worker = 10
	sshExecAgent.TimeOut = 30 * time.Second
	_, err := sshExecAgent.SshHostByKey(hosts, user, cmd)
	return err
}
