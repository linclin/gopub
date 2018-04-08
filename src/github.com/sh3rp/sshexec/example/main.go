package main

import (
	"fmt"
	"time"

	"github.com/sh3rp/sshexec"
)

func main() {
	agent := sshexec.NewAgent()
	agent.AddListener(func(r *sshexec.ExecResult) {
		fmt.Printf("%s: %s (%dms)\n", r.Host, r.Command, (r.EndTime.UnixNano()-r.StartTime.UnixNano())/1000)
	})
	agent.Start()
	agent.RunWithCreds("user", "password", "localhost", 22, "uname -a")
	agent.RunWithCreds("user", "password", "localhost", 22, "id")
	agent.RunWithCreds("user", "password", "localhost", 22, "ls -al")
	// spin

	for {
		time.Sleep(1000)
	}
}
