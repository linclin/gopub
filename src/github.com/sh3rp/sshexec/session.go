package sshexec

import (
	"bytes"
	"github.com/pkg/sftp"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/ssh"
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
	Id             uuid.UUID
	Host           string
	Command        string
	LocalFilePath  string
	RemoteFilePath string
	Result         string
	StartTime      time.Time
	EndTime        time.Time
	Error          error
}

// execute the command and return a result structure

func (exec *HostSession) Exec(id uuid.UUID, command string, config ssh.ClientConfig) *ExecResult {

	result := &ExecResult{
		Id:      id,
		Host:    exec.Hostname,
		Command: command,
	}

	client, err := ssh.Dial("tcp", exec.Hostname+":"+strconv.Itoa(exec.Port), &config)

	if err != nil {
		result.Error = err
		return result
	}

	session, err := client.NewSession()

	if err != nil {
		result.Error = err
		return result
	}

	defer session.Close()

	var b bytes.Buffer

	session.Stdout = &b
	var b1 bytes.Buffer
	session.Stderr = &b1
	start := time.Now()
	if err := session.Run(command); err != nil {
		result.Error = err
		result.Result = b1.String()
		return result
	}
	end := time.Now()
	result.Result = b.String()
	result.StartTime = start
	result.EndTime = end
	return result
}

// execute the command and return a result structure

func (exec *HostSession) Transfer(id uuid.UUID, localFilePath string, remoteFilePath string, config ssh.ClientConfig) *ExecResult {

	result := &ExecResult{
		Id:             id,
		Host:           exec.Hostname,
		LocalFilePath:  localFilePath,
		RemoteFilePath: remoteFilePath,
	}

	client, err := ssh.Dial("tcp", exec.Hostname+":"+strconv.Itoa(exec.Port), &config)

	if err != nil {
		result.Error = err
		return result
	}

	session, err := client.NewSession()

	if err != nil {
		result.Error = err
		return result
	}

	defer session.Close()

	srcFile, err := os.Open(localFilePath)
	if err != nil {
		result.Error = err
		return result
	}
	defer srcFile.Close()

	sftpClient, err := sftp.NewClient(client)
	// 这里换成实际的 SSH 连接的 用户名，密码，主机名或IP，SSH端口
	// create sftp client
	if err != nil {
		result.Error = err
		return result
	}
	defer sftpClient.Close()

	dstFile, err := sftpClient.Create(remoteFilePath)
	if err != nil {
		result.Error = err
		return result
	}
	defer dstFile.Close()
	buf := make([]byte, 1024)
	start := time.Now()
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		_, err := dstFile.Write(buf)
		if err != nil {
			result.Error = err
			return result
		}
	}
	end := time.Now()
	result.StartTime = start
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

	config.Ciphers = []string{"aes128-cbc", "3des-cbc"}

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
