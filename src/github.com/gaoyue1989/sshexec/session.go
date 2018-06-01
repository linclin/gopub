package sshexec

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

// ssh session

type HostSession struct {
	Username string
	Password string
	Hostname string
	Signers  []ssh.Signer
	Port     int
	Auths    []ssh.AuthMethod
}

// result of the command execution

type ExecResult struct {
	Id             int
	Host           string
	Command        string
	LocalFilePath  string
	RemoteFilePath string
	Result         string
	StartTime      time.Time
	EndTime        time.Time
	Error          error
	ErrorInfo          string
}

// execute the command and return a result structure

func (exec *HostSession) Exec(id int, command string, config ssh.ClientConfig) *ExecResult {

	result := &ExecResult{
		Id:      id,
		Host:    exec.Hostname,
		Command: command,
	}

	client, err := ssh.Dial("tcp", exec.Hostname+":"+strconv.Itoa(exec.Port), &config)

	if err != nil {
		result.Error = err
		result.ErrorInfo=err.Error()
		result.Result = "ssh连接出错，请检查key信任或ip是否正确"
		return result
	}

	session, err := client.NewSession()

	if err != nil {
		result.ErrorInfo=err.Error()
		result.Result = "ssh连接出错，请检查key信任或ip是否正确"
		result.Error = err
		return result
	}

	defer func() {
		session.Close()
		client.Close()
	}()

	var b bytes.Buffer

	session.Stdout = &b
	var b1 bytes.Buffer
	session.Stderr = &b1
	start := time.Now()
	if err := session.Run(command); err != nil {
		result.Error = err
		result.ErrorInfo=err.Error()
		result.Result = b1.String() +" 错误信息："+err.Error()
		return result
	}
	end := time.Now()
	result.Result = b.String()
	result.StartTime = start
	result.EndTime = end
	return result
}

// execute the command and return a result structure

func (exec *HostSession) Transfer(id int, localFilePath string, remoteFilePath string, config ssh.ClientConfig) *ExecResult {

	result := &ExecResult{
		Id:             id,
		Host:           exec.Hostname,
		LocalFilePath:  localFilePath,
		RemoteFilePath: remoteFilePath,
	}
	start := time.Now()
	result.StartTime = start
	client, err := ssh.Dial("tcp", exec.Hostname+":"+strconv.Itoa(exec.Port), &config)

	if err != nil {
		result.ErrorInfo=err.Error()
		result.Result = "ssh连接出错，请检查key信任或ip是否正确"
		result.Error = err
		return result
	}

	session, err := client.NewSession()

	if err != nil {
		result.ErrorInfo=err.Error()
		result.Result = "ssh连接出错，请检查key信任或ip是否正确"
		result.Error = err
		return result
	}

	defer func() {
		session.Close()
		client.Close()
	}()
	var fileSize int64
	if s, err := os.Stat(localFilePath); err != nil {
		result.Error = err
		return result

	} else {
		fileSize = s.Size()
	}
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		result.ErrorInfo=err.Error()
		result.Result = err.Error()
		result.Error = err
		return result
	}

	defer srcFile.Close()

	sftpClient, err := sftp.NewClient(client)
	// 这里换成实际的 SSH 连接的 用户名，密码，主机名或IP，SSH端口
	// create sftp client
	if err != nil {
		result.ErrorInfo=err.Error()
		result.Result = err.Error()
		result.Error = err
		return result
	}
	defer sftpClient.Close()

	dstFile, err := sftpClient.Create(remoteFilePath)
	if err != nil {
		result.ErrorInfo=err.Error()
		result.Result = err.Error()
		result.Error = err
		return result
	}
	defer dstFile.Close()
	//todo 这里使用buff池 或io.Copy 性能没有测试出差距
	//n, err := Copy(dstFile, io.LimitReader(srcFile, fileSize))
	n, err := io.Copy(dstFile, io.LimitReader(srcFile, fileSize))
	if err != nil {
		result.ErrorInfo=err.Error()
		result.Result = err.Error()
		result.Error = err
		return result
	}
	if n != fileSize {
		beego.Info(err)
		result.Error = errors.New(fmt.Sprintf("copy: expected %v bytes, got %d", fileSize, n))
		result.ErrorInfo=result.Error.Error()
		result.Result = result.Error.Error()
		return result
	}
	end := time.Now()
	result.EndTime = end
	return result
}

func (exec *HostSession) GenerateConfig() ssh.ClientConfig {
	var auths []ssh.AuthMethod

	if len(exec.Password) != 0 {
		auths = append(auths, ssh.Password(exec.Password))
	} else {
		if len(exec.Auths) > 0 {
			auths = exec.Auths
		} else {
			auths = append(auths, ssh.PublicKeys(exec.Signers...))
		}

	}

	config := ssh.ClientConfig{
		User: exec.Username,
		Auth: auths,
	}

	config.Ciphers = []string{"aes128-cbc", "3des-cbc","aes128-ctr","aes192-ctr","aes256-ctr" }

	return config
}

func readKey(filename string) (ssh.Signer, error) {
	f, _ := os.Open(filename)

	bytes, _ := ioutil.ReadAll(f)
	return generateKey(bytes)
}

func generateKey(keyData []byte) (ssh.Signer, error) {
	return ssh.ParsePrivateKey(keyData)
}
