package common

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// HTTPGet return the body of the response when send http get method to the server
func (s *BaseService) HTTPGet(addr, urlpath string) (rspBody []byte, err error) {
	return SendHTTPReq(s.Cfg, "GET", addr, urlpath, nil)
}

// HTTPPost return the body of the response when send http get method to the server
func (s *BaseService) HTTPPost(addr, urlpath string, reqBody []byte) (rspBody []byte, err error) {
	return SendHTTPReq(s.Cfg, "POST", addr, urlpath, reqBody)
}

// HTTPDelete return the body of the response when send http get method to the server
func (s *BaseService) HTTPDelete(addr, urlpath string) (err error) {
	_, err = SendHTTPReq(s.Cfg, "DELETE", addr, urlpath, nil)
	return
}

// CreateHTTPClient return a http client instannce
func CreateHTTPClient(cfg *Config) *http.Client {
	var client *http.Client
	tr := &http.Transport{
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   1,
		DisableKeepAlives:     true,
	}

	client = &http.Client{Transport: tr}
	return client
}

// SendHTTPReqWithClient ...
func SendHTTPReqWithClient(client *http.Client, cfg *Config, method, addr, urlpath string, reqBody []byte) (rspBody []byte, err error) {
	schema := "http"

	if cfg.Server && !strings.Contains(addr, ":") {
		addr = fmt.Sprintf("%s:%v", addr, cfg.Net.AgentMgntPort)
	}

	url := fmt.Sprintf("%s://%s%s", schema, addr, urlpath)
	req, err := http.NewRequest(method, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(cfg.Auth.Username, cfg.Auth.Password)
	req.Header.Set("Content-Type", "application/json")
	//log.Debugf("Sending http request %v", req)

	client.Timeout = 2 * time.Second
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 300 {
		return nil, fmt.Errorf("Recv http status code %v", resp.StatusCode)
	}

	if resp.ContentLength > 0 {
		rspBody, err = ioutil.ReadAll(resp.Body)
	}
	return
}

// SendHTTPReq ...
func SendHTTPReq(cfg *Config, method, addr, urlpath string, reqBody []byte) (rspBody []byte, err error) {
	client := CreateHTTPClient(cfg)
	return SendHTTPReqWithClient(client, cfg, method, addr, urlpath, reqBody)
}
