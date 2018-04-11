// Copyright 2015 Vadim Kravcenko
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Gojenkins is a Jenkins Client in Go, that exposes the jenkins REST api in a more developer friendly way.
package gojenkins

import (
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// Basic Authentication
type BasicAuth struct {
	Username string
	Password string
}

type Jenkins struct {
	Server    string
	Version   string
	Raw       *executorResponse
	Requester *Requester
}

// Loggers
var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

// Init Method. Should be called after creating a Jenkins Instance.
// e.g jenkins := CreateJenkins("url").Init()
// HTTP Client is set here, Connection to jenkins is tested here.
func (j *Jenkins) Init() (*Jenkins, error) {
	j.initLoggers()
	// Skip SSL Verification?
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !j.Requester.SslVerify},
	}

	if j.Requester.Client == nil {
		cookies, _ := cookiejar.New(nil)

		if os.Getenv("HTTP_PROXY") != "" {
			proxyUrl, _ := url.Parse(os.Getenv("HTTP_PROXY"))
			tr.Proxy = http.ProxyURL(proxyUrl)
		}

		client := &http.Client{
			Transport: tr,
			Jar:       cookies,
			// Fucntion to add auth on redirect.
			CheckRedirect: j.Requester.redirectPolicyFunc,
		}

		j.Requester.Client = client
	}

	// Check Connection
	j.Raw = new(executorResponse)
	rsp, err := j.Requester.GetJSON("/", j.Raw, nil)

	if err != nil {
		return nil, err
	}

	j.Version = rsp.Header.Get("X-Jenkins")
	if j.Raw == nil {
		return nil, errors.New("Connection Failed, Please verify that the host and credentials are correct.")
	}

	return j, nil
}

func (j *Jenkins) initLoggers() {
	Info = log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(os.Stdout,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(os.Stderr,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

// Get Basic Information About Jenkins
func (j *Jenkins) Info() (*executorResponse, error) {
	_, err := j.Requester.Get("/", j.Raw, nil)

	if err != nil {
		return nil, err
	}
	return j.Raw, nil
}

// Create a new Node
// Can be JNLPLauncher or SSHLauncher
// Example : jenkins.CreateNode("nodeName", 1, "Description", "/var/lib/jenkins", map[string]string{"method": "JNLPLauncher"})
// By Default JNLPLauncher is created
func (j *Jenkins) CreateNode(name string, numExecutors int, description string, remoteFS string, options ...interface{}) (*Node, error) {
	node, _ := j.GetNode(name)

	if node != nil {
		return node, nil
	}

	params := map[string]string{"method": "JNLPLauncher"}

	if len(options) > 0 {
		params, _ = options[0].(map[string]string)
	}

	if _, ok := params["method"]; !ok {
		params["method"] = "JNLPLauncher"
	}

	method := params["method"]
	var launcher map[string]string
	switch method {
	case "":
		fallthrough
	case "JNLPLauncher":
		launcher = map[string]string{"stapler-class": "hudson.slaves.JNLPLauncher"}
	case "SSHLauncher":
		launcher = map[string]string{
			"stapler-class":        "hudson.plugins.sshslaves.SSHLauncher",
			"$class":               "hudson.plugins.sshslaves.SSHLauncher",
			"host":                 params["host"],
			"port":                 params["port"],
			"credentialsId":        params["credentialsId"],
			"jvmOptions":           params["jvmOptions"],
			"javaPath":             params["javaPath"],
			"prefixStartSlaveCmd":  params["prefixStartSlaveCmd"],
			"suffixStartSlaveCmd":  params["suffixStartSlaveCmd"],
			"maxNumRetries":        params["maxNumRetries"],
			"retryWaitTime":        params["retryWaitTime"],
			"lanuchTimeoutSeconds": params["lanuchTimeoutSeconds"],
			"type":                 "hudson.slaves.DumbSlave",
			"stapler-class-bag":    "true"}
	default:
		return nil, errors.New("launcher method not supported")
	}

	node = &Node{Jenkins: j, Raw: new(NodeResponse), Base: "/computer/" + name}
	NODE_TYPE := "hudson.slaves.DumbSlave$DescriptorImpl"
	MODE := "NORMAL"
	qr := map[string]string{
		"name": name,
		"type": NODE_TYPE,
		"json": makeJson(map[string]interface{}{
			"name":               name,
			"nodeDescription":    description,
			"remoteFS":           remoteFS,
			"numExecutors":       numExecutors,
			"mode":               MODE,
			"type":               NODE_TYPE,
			"retentionsStrategy": map[string]string{"stapler-class": "hudson.slaves.RetentionStrategy$Always"},
			"nodeProperties":     map[string]string{"stapler-class-bag": "true"},
			"launcher":           launcher,
		}),
	}

	resp, err := j.Requester.Post("/computer/doCreateItem", nil, nil, qr)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 400 {
		_, err := node.Poll()
		if err != nil {
			return nil, err
		}
		return node, nil
	}
	return nil, errors.New(strconv.Itoa(resp.StatusCode))
}

// Create a new job from config File
// Method takes XML string as first parameter, and if the name is not specified in the config file
// takes name as string as second parameter
// e.g jenkins.CreateJob("<config></config>","newJobName")
func (j *Jenkins) CreateJob(config string, options ...interface{}) (*Job, error) {
	qr := make(map[string]string)
	if len(options) > 0 {
		qr["name"] = options[0].(string)
	} else {
		return nil, errors.New("Error Creating Job, job name is missing")
	}
	jobObj := Job{Jenkins: j, Raw: new(jobResponse), Base: "/job/" + qr["name"]}
	job, err := jobObj.Create(config, qr)
	if err != nil {
		return nil, err
	}
	return job, nil
}

// Rename a job.
// First parameter job old name, Second parameter job new name.
func (j *Jenkins) RenameJob(job string, name string) *Job {
	jobObj := Job{Jenkins: j, Raw: new(jobResponse), Base: "/job/" + job}
	jobObj.Rename(name)
	return &jobObj
}

// Create a copy of a job.
// First parameter Name of the job to copy from, Second parameter new job name.
func (j *Jenkins) CopyJob(copyFrom string, newName string) (*Job, error) {
	job := Job{Jenkins: j, Raw: new(jobResponse), Base: "/job/" + copyFrom}
	_, err := job.Poll()
	if err != nil {
		return nil, err
	}
	return job.Copy(newName)
}

// Delete a job.
func (j *Jenkins) DeleteJob(name string) (bool, error) {
	job := Job{Jenkins: j, Raw: new(jobResponse), Base: "/job/" + name}
	return job.Delete()
}

// Invoke a job.
// First parameter job name, second parameter is optional Build parameters.
func (j *Jenkins) BuildJob(name string, options ...interface{}) (int64, error) {
	job := Job{Jenkins: j, Raw: new(jobResponse), Base: "/job/" + name}
	var params map[string]string
	if len(options) > 0 {
		params, _ = options[0].(map[string]string)
	}
	return job.InvokeSimple(params)
}

func (j *Jenkins) GetNode(name string) (*Node, error) {
	node := Node{Jenkins: j, Raw: new(NodeResponse), Base: "/computer/" + name}
	status, err := node.Poll()
	if err != nil {
		return nil, err
	}
	if status == 200 {
		return &node, nil
	}
	return nil, errors.New("No node found")
}

func (j *Jenkins) GetLabel(name string) (*Label, error) {
	label := Label{Jenkins: j, Raw: new(LabelResponse), Base: "/label/" + name}
	status, err := label.Poll()
	if err != nil {
		return nil, err
	}
	if status == 200 {
		return &label, nil
	}
	return nil, errors.New("No label found")
}

func (j *Jenkins) GetBuild(jobName string, number int64) (*Build, error) {
	job, err := j.GetJob(jobName)
	if err != nil {
		return nil, err
	}
	build, err := job.GetBuild(number)

	if err != nil {
		return nil, err
	}
	return build, nil
}

func (j *Jenkins) GetJob(id string, parentIDs ...string) (*Job, error) {
	job := Job{Jenkins: j, Raw: new(jobResponse), Base: "/job/" + strings.Join(append(parentIDs, id), "/job/")}
	status, err := job.Poll()
	if err != nil {
		return nil, err
	}
	if status == 200 {
		return &job, nil
	}
	return nil, errors.New(strconv.Itoa(status))
}

func (j *Jenkins) GetAllNodes() ([]*Node, error) {
	computers := new(Computers)

	qr := map[string]string{
		"depth": "1",
	}

	_, err := j.Requester.GetJSON("/computer", computers, qr)
	if err != nil {
		return nil, err
	}

	nodes := make([]*Node, len(computers.Computers))
	for i, node := range computers.Computers {
		nodes[i] = &Node{Jenkins: j, Raw: node, Base: "/computer/" + node.DisplayName}
	}

	return nodes, nil
}

// Get all builds Numbers and URLS for a specific job.
// There are only build IDs here,
// To get all the other info of the build use jenkins.GetBuild(job,buildNumber)
// or job.GetBuild(buildNumber)
func (j *Jenkins) GetAllBuildIds(job string) ([]jobBuild, error) {
	jobObj, err := j.GetJob(job)
	if err != nil {
		return nil, err
	}
	return jobObj.GetAllBuildIds()
}

// Get Only Array of Job Names, Color, URL
// Does not query each single Job.
func (j *Jenkins) GetAllJobNames() ([]job, error) {
	exec := Executor{Raw: new(executorResponse), Jenkins: j}
	_, err := j.Requester.GetJSON("/", exec.Raw, nil)

	if err != nil {
		return nil, err
	}

	return exec.Raw.Jobs, nil
}

// Get All Possible Job Objects.
// Each job will be queried.
func (j *Jenkins) GetAllJobs() ([]*Job, error) {
	exec := Executor{Raw: new(executorResponse), Jenkins: j}
	_, err := j.Requester.GetJSON("/", exec.Raw, nil)

	if err != nil {
		return nil, err
	}

	jobs := make([]*Job, len(exec.Raw.Jobs))
	for i, job := range exec.Raw.Jobs {
		ji, err := j.GetJob(job.Name)
		if err != nil {
			return nil, err
		}
		jobs[i] = ji
	}
	return jobs, nil
}

// Returns a Queue
func (j *Jenkins) GetQueue() (*Queue, error) {
	q := &Queue{Jenkins: j, Raw: new(queueResponse), Base: j.GetQueueUrl()}
	_, err := q.Poll()
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (j *Jenkins) GetQueueUrl() string {
	return "/queue"
}

// Get Artifact data by Hash
func (j *Jenkins) GetArtifactData(id string) (*fingerPrintResponse, error) {
	fp := Fingerprint{Jenkins: j, Base: "/fingerprint/", Id: id, Raw: new(fingerPrintResponse)}
	return fp.GetInfo()
}

// Returns the list of all plugins installed on the Jenkins server.
// You can supply depth parameter, to limit how much data is returned.
func (j *Jenkins) GetPlugins(depth int) (*Plugins, error) {
	p := Plugins{Jenkins: j, Raw: new(pluginResponse), Base: "/pluginManager", Depth: depth}
	_, err := p.Poll()
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Check if the plugin is installed on the server.
// Depth level 1 is used. If you need to go deeper, you can use GetPlugins, and iterate through them.
func (j *Jenkins) HasPlugin(name string) (*Plugin, error) {
	p, err := j.GetPlugins(1)

	if err != nil {
		return nil, err
	}
	return p.Contains(name), nil
}

// Verify Fingerprint
func (j *Jenkins) ValidateFingerPrint(id string) (bool, error) {
	fp := Fingerprint{Jenkins: j, Base: "/fingerprint/", Id: id, Raw: new(fingerPrintResponse)}
	valid, err := fp.Valid()
	if err != nil {
		return false, err
	}
	if valid {
		return true, nil
	}
	return false, nil
}

func (j *Jenkins) GetView(name string) (*View, error) {
	url := "/view/" + name
	view := View{Jenkins: j, Raw: new(viewResponse), Base: url}
	_, err := view.Poll()
	if err != nil {
		return nil, err
	}
	return &view, nil
}

func (j *Jenkins) GetAllViews() ([]*View, error) {
	_, err := j.Poll()
	if err != nil {
		return nil, err
	}
	views := make([]*View, len(j.Raw.Views))
	for i, v := range j.Raw.Views {
		views[i], _ = j.GetView(v.Name)
	}
	return views, nil
}

// Create View
// First Parameter - name of the View
// Second parameter - Type
// Possible Types:
// 		gojenkins.LIST_VIEW
// 		gojenkins.NESTED_VIEW
// 		gojenkins.MY_VIEW
// 		gojenkins.DASHBOARD_VIEW
// 		gojenkins.PIPELINE_VIEW
// Example: jenkins.CreateView("newView",gojenkins.LIST_VIEW)
func (j *Jenkins) CreateView(name string, viewType string) (*View, error) {
	exists, err := j.GetView(name)
	if err != nil {
		return nil, err
	}
	if exists.Raw.Name != "" {
		Error.Println("View Already exists.")
		return exists, errors.New("View already exists")
	}
	view := &View{Jenkins: j, Raw: new(viewResponse), Base: "/view/" + name}
	endpoint := "/createView"
	data := map[string]string{
		"name":   name,
		"mode":   viewType,
		"Submit": "OK",
		"json": makeJson(map[string]string{
			"name": name,
			"mode": viewType,
		}),
	}
	r, err := j.Requester.Post(endpoint, nil, view.Raw, data)

	if err != nil {
		return nil, err
	}

	if r.StatusCode == 200 {
		return j.GetView(name)
	}
	return nil, errors.New(strconv.Itoa(r.StatusCode))
}

func (j *Jenkins) Poll() (int, error) {
	resp, err := j.Requester.GetJSON("/", j.Raw, nil)
	if err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}

// Creates a new Jenkins Instance
// Optional parameters are: username, password
// After creating an instance call init method.
func CreateJenkins(base string, auth ...interface{}) *Jenkins {
	j := &Jenkins{}
	if strings.HasSuffix(base, "/") {
		base = base[:len(base)-1]
	}
	j.Server = base
	j.Requester = &Requester{Base: base, SslVerify: false}
	if len(auth) == 2 {
		j.Requester.BasicAuth = &BasicAuth{Username: auth[0].(string), Password: auth[1].(string)}
	}
	return j
}
