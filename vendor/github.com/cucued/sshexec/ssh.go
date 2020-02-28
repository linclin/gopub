package sshexec

import (
	"time"
	"io/ioutil"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"github.com/linclin/grpool"
	"errors"
)

//
// Main agent struct
//

type SSHExecAgent struct {
	Worker  int
	TimeOut time.Duration
}

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func GetAuthPassword(password string) []ssh.AuthMethod {
	return []ssh.AuthMethod{ssh.Password(password)}
}

func GetAuthKeys(keys []string) []ssh.AuthMethod {
	methods := []ssh.AuthMethod{}
	for _, keyname := range keys {
		pkey := PublicKeyFile(keyname)
		if pkey != nil {
			methods = append(methods, pkey)
		}
	}
	return methods
}
func (s *SSHExecAgent) SshHostByKey(hosts []string, port int, user string, cmd string) ([]ExecResult, error) {
	if (len(hosts) == 0) {
		log.Println("no hosts")
		return nil, errors.New("no hosts")
	}
	if (s.Worker == 0) {
		s.Worker = 40
	}
	if (s.TimeOut == 0) {
		s.TimeOut = 3600 * time.Second
	}
	keys := []string{
		os.Getenv("HOME") + "/.ssh/id_dsa",
		os.Getenv("HOME") + "/.ssh/id_rsa",
	}
	authKeys := GetAuthKeys(keys)
	if len(authKeys) < 1 {
		log.Println("the user no key")
		return nil, errors.New("the user no key")
	}
	pool := grpool.NewPool(s.Worker, len(hosts), s.TimeOut)
	defer pool.Release()
	pool.WaitCount(len(hosts))
	for i, _ := range hosts {
		count := i
		pool.JobQueue <- grpool.Job{
			Jobid: count,
			Jobfunc: func() (interface{}, error) {
				session := &HostSession{
					Username: user,
					Password: "",
					Hostname: hosts[count],
					Port:     port,
					Auths:  authKeys,
				}
				r := session.Exec(count, cmd, session.GenerateConfig())
				return *r, nil
			},
		}
	}

	pool.WaitAll()
	returnResult := make([]ExecResult, len(hosts))
	errorText := ""
	for res := range pool.Jobresult {
		jobId, _ := res.Jobid.(int)
		if (res.Timedout) {
			returnResult[jobId].Id = jobId
			returnResult[jobId].Host = hosts[jobId]
			returnResult[jobId].Command = cmd
			returnResult[jobId].Error = errors.New("ssh time out")
			errorText += "the host " + hosts[jobId] + " commond  exec time out."
		} else {
			execResult, _ := res.Result.(ExecResult)
			returnResult[jobId] = execResult
			if (execResult.Error != nil) {
				errorText += "the host " + execResult.Host + " commond  exec error.\n" + "rsult info :" + execResult.Result + ".\nerror info :" + execResult.Error.Error()
			}
		}
	}
	if (errorText != "") {
		return returnResult, errors.New(errorText)

	} else {
		return returnResult, nil
	}

}

func (s *SSHExecAgent) SftpHostByKey(hosts []string, port int, user string, localFilePath  string, remoteFilePath string) ([]ExecResult, error) {
	if (len(hosts) == 0) {
		log.Println("no hosts")
		return nil, errors.New("no hosts")
	}
	if (s.Worker == 0) {
		s.Worker = 40
	}
	if (s.TimeOut == 0) {
		s.TimeOut = 3600 * time.Second
	}
	keys := []string{
		os.Getenv("HOME") + "/.ssh/id_dsa",
		os.Getenv("HOME") + "/.ssh/id_rsa",
	}
	authKeys := GetAuthKeys(keys)
	if len(authKeys) < 1 {
		log.Println("the user no key")
		return nil, errors.New("the user no key")
	}
	pool := grpool.NewPool(s.Worker, len(hosts), s.TimeOut)
	defer pool.Release()
	pool.WaitCount(len(hosts))
	for i, _ := range hosts {
		count := i
		pool.JobQueue <- grpool.Job{
			Jobid: count,
			Jobfunc: func() (interface{}, error) {
				session := &HostSession{
					Username: user,
					Password: "",
					Hostname: hosts[count],
					Port:     port,
					Auths:  authKeys,
				}
				r := session.Transfer(count, localFilePath, remoteFilePath, session.GenerateConfig())
				return *r, nil
			},
		}
	}

	pool.WaitAll()
	returnResult := make([]ExecResult, len(hosts))
	errorText := ""
	for res := range pool.Jobresult {
		jobId, _ := res.Jobid.(int)
		if (res.Timedout) {
			returnResult[jobId].Id = jobId
			returnResult[jobId].Host = hosts[jobId]
			returnResult[jobId].LocalFilePath = localFilePath
			returnResult[jobId].RemoteFilePath = remoteFilePath
			returnResult[jobId].Error = errors.New("ssh time out")
			errorText += "the host " + hosts[jobId] + " commond  exec time out."
		} else {
			execResult, _ := res.Result.(ExecResult)
			returnResult[jobId] = execResult
			if (execResult.Error != nil) {
				errorText += "the host " + execResult.Host + " commond  exec error.\n" + "rsult info :" + execResult.Result + ".\nerror info :" + execResult.Error.Error()
			}
		}
	}
	if (errorText != "") {
		return returnResult, errors.New(errorText)

	} else {
		return returnResult, nil
	}

}
