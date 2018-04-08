package server

import (
	"github.com/julienschmidt/httprouter"
	"library/p2p/common"
	"library/p2p/p2p"
	"time"
)

type Server struct {
	common.BaseService
	// 用于缓存当前接收到任务
	cache *common.Cache
	// Session管理
	sessionMgnt *p2p.TaskSessionMgnt
}

func NewServer(cfg *common.Config) (*Server, error) {
	s := &Server{
		cache:       common.NewCache(5 * time.Minute),
		sessionMgnt: p2p.NewSessionMgnt(cfg),
	}
	s.BaseService = *common.NewBaseService(cfg, cfg.Name, s)
	return s, nil
}

func (s *Server) OnStart(cfg *common.Config, e *httprouter.Router) error {
	go func() { s.sessionMgnt.Start() }()
	e.POST("/api/v1/server/tasks", s.CreateTask)
	e.DELETE("/api/v1/server/tasks/:id", s.CancelTask)
	e.GET("/api/v1/server/tasks/:id", s.QueryTask)
	e.POST("/api/v1/server/tasks/status", s.ReportTask)
	return nil
}

func (s *Server) OnStop(c *common.Config, e *httprouter.Router) {
	go func() { s.sessionMgnt.Stop() }()
}
