package sshexec

import (
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/ssh"
)

//
// Main agent struct
//

type SSHExecAgent struct {
	results   chan *ExecResult
	AllReturn []ExecResult
	listeners []func(*ExecResult)
	running   bool
}

// constructor

func NewAgent() *SSHExecAgent {
	agent := &SSHExecAgent{
		results: make(chan *ExecResult, 100),
		running: false,
	}
	agent.Start()
	return agent
}

//
// Runs a command with specified credentials (username/password)
//

func (agent *SSHExecAgent) RunWithCreds(username string, password string, authMethod []ssh.AuthMethod, hostname string, port int, command string) *uuid.UUID {
	session := &HostSession{
		Username: username,
		Password: password,
		Hostname: hostname,
		Port:     port,
		Auths:    authMethod,
	}

	return agent.RunWithSession(session, command)
}

func (agent *SSHExecAgent) TransferWithCreds(username string, password string, authMethod []ssh.AuthMethod, hostname string, port int, localFilePath string, remoteFilePath string) *uuid.UUID {
	session := &HostSession{
		Username: username,
		Password: password,
		Hostname: hostname,
		Port:     port,
		Auths:    authMethod,
	}

	return agent.TransferWithSession(session, localFilePath, remoteFilePath)
}

//
// Runs a command with a specified session
//

func (agent *SSHExecAgent) RunWithSession(session *HostSession, command string) *uuid.UUID {
	id := uuid.NewV4()
	go func(uuid uuid.UUID) {
		r := session.Exec(uuid, command, session.GenerateConfig())
		agent.results <- r
	}(id)
	return &id
}

//
// Runs a command with a specified session
//

func (agent *SSHExecAgent) TransferWithSession(session *HostSession, localFilePath string, remoteFilePath string) *uuid.UUID {
	id := uuid.NewV4()
	go func(uuid uuid.UUID) {
		r := session.Transfer(uuid, localFilePath, remoteFilePath, session.GenerateConfig())
		agent.results <- r
	}(id)
	return &id
}

//
// Add an ExecResult listener
//

func (agent *SSHExecAgent) AddListener(f func(*ExecResult)) {
	agent.listeners = append(agent.listeners, f)
}

//
// Start the agent result channel and publish results as they come in to the channel
//

func (agent *SSHExecAgent) Start() {
	if agent.running {
		return
	}

	agent.running = true
	go func() {
		for agent.running {
			select {
			case result := <-agent.results:
				if len(agent.listeners) > 0 {
					for _, listener := range agent.listeners {
						listener(result)
					}
				}
			}
		}
	}()
}

//
// Stop the agent results channel
//

func (agent *SSHExecAgent) Stop() {
	agent.running = false
}
