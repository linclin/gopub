package server

import (
	"net/http"

	"encoding/json"
	"errors"
	log "github.com/cihub/seelog"
	"github.com/julienschmidt/httprouter"
	"github.com/linclin/gopub/src/library/p2p/p2p"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

const (
	// For use with functions that take an expiration time.
	NoExpiration time.Duration = -1
)

func (svc *Server) String(r int, s string, w http.ResponseWriter) {
	w.WriteHeader(r)
	w.Write([]byte(s))
}
func (svc *Server) Json(r int, s interface{}, w http.ResponseWriter) {
	w.WriteHeader(r)
	ss, _ := json.Marshal(s)
	w.Write(ss)
}
func (svc *Server) getRequestParams(r *http.Request, s interface{}) error {
	if r.Body == nil {
		return nil
	}
	defer r.Body.Close()
	rbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	} else {
		if err := json.Unmarshal(rbody, &s); err != nil {
			return err
		}
	}
	return nil
}

// CreateTask POST /api/v1/server/tasks
func (s *Server) CreateTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//  获取Body
	t := new(CreateTask)
	if err := s.getRequestParams(r, t); err != nil {
		log.Errorf("Recv [%s] request, decode body failed. %v", "/api/v1/server/tasks", err)
		return
	}

	// 检查任务是否存在
	v, ok := s.cache.Get(t.ID)
	if ok {
		cti := v.(*CachedTaskInfo)
		if cti.EqualCmp(t) {
			s.String(http.StatusAccepted, "", w)
			return
		}
		log.Debugf("[%s] Recv task, task is existed", t.ID)
		s.String(http.StatusBadRequest, TaskExist.String(), w)
		return
	}

	log.Infof("[%s] Recv task, file=%v, ips=%v", t.ID, t.DispatchFiles, t.DestIPs)

	cti := NewCachedTaskInfo(s, t)
	s.cache.Set(t.ID, cti, NoExpiration)
	s.cache.OnEvicted(func(id string, v interface{}) {
		log.Infof("[%s] Remove task cache", t.ID)
		cti := v.(*CachedTaskInfo)
		cti.quitChan <- struct{}{}
	})
	go cti.Start()

	s.String(http.StatusAccepted, "", w)
	return
}

// CancelTask DELETE /api/v1/server/tasks/:id
func (s *Server) CancelTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	log.Infof("[%s] Recv cancel task", id)
	v, ok := s.cache.Get(id)
	if !ok {
		s.String(http.StatusBadRequest, TaskNotExist.String(), w)
		return
	}
	cti := v.(*CachedTaskInfo)
	cti.stopChan <- struct{}{}
	s.Json(http.StatusAccepted, "", w)
	return
}

// QueryTask GET /api/v1/server/tasks/:id
func (s *Server) QueryTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	log.Infof("[%s] Recv query task", id)
	v, ok := s.cache.Get(id)
	if !ok {
		s.String(http.StatusBadRequest, TaskNotExist.String(), w)
		return
	}
	cti := v.(*CachedTaskInfo)
	s.Json(http.StatusOK, cti.Query(), w)
	return
}

// ReportTask POST /api/v1/server/tasks/status
func (s *Server) ReportTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//  获取Body
	csr := new(p2p.StatusReport)
	if err := s.getRequestParams(r, csr); err != nil {
		log.Errorf("Recv [%s] request, decode body failed. %v", "", err)
		return
	}

	log.Debugf("[%s] Recv task report, ip=%v, pecent=%v", csr.TaskID, csr.IP, csr.PercentComplete)
	if v, ok := s.cache.Get(csr.TaskID); ok {
		cti := v.(*CachedTaskInfo)
		cti.reportChan <- csr
	}

	s.String(http.StatusOK, "", w)
	return
}

func (s *Server) QueryTaskNoHttp(id string) (*TaskInfo, error) {
	log.Infof("[%s] Recv query task", id)
	v, ok := s.cache.Get(id)
	if !ok {
		return new(TaskInfo), errors.New(TaskNotExist.String())
	}
	cti := v.(*CachedTaskInfo)
	return cti.Query(), nil
}
func (s *Server) CreateTaskNoHttp(t *CreateTask) error {
	// 检查任务是否存在
	v, ok := s.cache.Get(t.ID)
	if ok {
		cti := v.(*CachedTaskInfo)
		if cti.EqualCmp(t) {
			return nil
		}
		log.Debugf("[%s] Recv task, task is existed", t.ID)

		return errors.New(TaskExist.String())
	}

	log.Infof("[%s] Recv task, file=%v, ips=%v", t.ID, t.DispatchFiles, t.DestIPs)

	cti := NewCachedTaskInfo(s, t)
	s.cache.Set(t.ID, cti, NoExpiration)
	s.cache.OnEvicted(func(id string, v interface{}) {
		log.Infof("[%s] Remove task cache", t.ID)
		cti := v.(*CachedTaskInfo)
		cti.quitChan <- struct{}{}
	})
	go cti.Start()
	return nil
}

// 给所有客户端发送停止命令
func (s *Server) CheckAllClientIp(ips []string) {
	url := "/api/v1/agent/ip/"
	for _, ip := range ips {
		go func(ip string) {
			sendip := ip
			if idx := strings.Index(ip, ":"); idx > 0 {
				ip = ip[:idx]
			}
			ip = ip + ":" + strconv.Itoa(s.Cfg.Net.AgentMgntPort)
			if _, err2 := s.HTTPGet(ip, url+sendip); err2 != nil {
				log.Errorf("Send http request failed. GET, ip=%s, url=%s, error=%v", ip, url, err2)
			} else {
				log.Debugf("Send http request success. GET, ip=%s, url=%s", ip, url)
			}
		}(ip)
	}
}

// 给所有客户端是否存在
func (s *Server) CheckAllClient(ips []string) map[string]string {
	res := map[string]string{}
	url := "/api/v1/agent/alive"
	for _, ip := range ips {
		if idx := strings.Index(ip, ":"); idx > 0 {
			ip = ip[:idx]
		}
		ip = ip + ":" + strconv.Itoa(s.Cfg.Net.AgentMgntPort)
		if _, err2 := s.HTTPGet(ip, url); err2 != nil {
			res[ip] = "dead"
		} else {
			res[ip] = "alive"
		}

	}
	return res
}
