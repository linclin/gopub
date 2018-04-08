package common

import (
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

type Service interface {
	Start() error
	Stop() bool
	OnStart(c *Config, e *httprouter.Router) error
	OnStop(c *Config, e *httprouter.Router)
	IsRunning() bool
}

type BaseService struct {
	name       string
	running    uint32 // atomic
	Cfg        *Config
	httprouter *httprouter.Router
	svc        Service
}

func NewBaseService(cfg *Config, name string, svc Service) *BaseService {
	return &BaseService{
		name:       name,
		running:    0,
		Cfg:        cfg,
		httprouter: httprouter.New(),
		svc:        svc,
	}
}

// init log by config
func (s *BaseService) initlog() {
	//以后添加日志模块

}

func (s *BaseService) runHttprouter() error {
	net := s.Cfg.Net
	log.Printf("Starting http server %s:%v", net.IP, net.MgntPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", net.MgntPort), s.httprouter)
	if err != nil {
		log.Printf("Start http server %s:%v failed %v", net.IP, net.MgntPort, err)
		return err
	}
	return nil
}

func (s *BaseService) Start() error {
	if atomic.CompareAndSwapUint32(&s.running, 0, 1) {
		s.initlog()
		log.Printf("Starting %s", s.name)
		if err := s.svc.OnStart(s.Cfg, s.httprouter); err != nil {
			return err
		}
		go func() {
			err := s.runHttprouter()
			if err != nil && !s.Cfg.Server {
				os.Exit(1)
			}
		}()
		return nil
	} else {
		return errors.New("Started aleadry.")
	}
}

func (s *BaseService) OnStart(c *Config, e *httprouter.Router) error { return nil }

func (s *BaseService) Stop() bool {
	if atomic.CompareAndSwapUint32(&s.running, 1, 0) {
		log.Printf("Stopping %s", s.name)
		s.svc.OnStop(s.Cfg, s.httprouter)
		return true
	} else {
		return false
	}
}

// Implements Service
func (s *BaseService) OnStop(c *Config, e *httprouter.Router) {}

// Implements Service
func (s *BaseService) IsRunning() bool {
	return atomic.LoadUint32(&s.running) == 1
}

func (s *BaseService) Auth(u, p string) bool {
	if u == s.Cfg.Auth.Username && p == s.Cfg.Auth.Password {
		return true
	}
	return false
}

// checks whether a file or directory exists.
// It returns false when the file or directory does not exist.
func FileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
