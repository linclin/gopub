package agent

import (
	"github.com/julienschmidt/httprouter"
	"library/p2p/common"
	"library/p2p/p2p"
	"os"
)

type Agent struct {
	common.BaseService
	// Session管理
	sessionMgnt *p2p.TaskSessionMgnt
}

func NewAgent(cfg *common.Config) (*Agent, error) {
	c := &Agent{
		sessionMgnt: p2p.NewSessionMgnt(cfg),
	}
	c.BaseService = *common.NewBaseService(cfg, cfg.Name, c)
	return c, nil
}

func (c *Agent) OnStart(cfg *common.Config, e *httprouter.Router) error {
	go func() {
		err := c.sessionMgnt.Start()
		if err != nil {
			os.Exit(1)
		}
	}()

	e.POST("/api/v1/agent/tasks", c.CreateTask)
	e.POST("/api/v1/agent/tasks/start", c.StartTask)
	e.DELETE("/api/v1/agent/tasks/:id", c.CancelTask)
	e.GET("/api/v1/agent/ip/:ip", c.ChangeIp)
	e.GET("/api/v1/agent/alive", c.Alive)
	return nil
}

func (c *Agent) OnStop(cfg *common.Config, e *httprouter.Router) {
	go func() { c.sessionMgnt.Stop() }()
}
