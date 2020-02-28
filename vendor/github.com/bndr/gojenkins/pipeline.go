// Copyright 2017 - Tessa Nordgren
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
//
// this file implements the pipeline-stage-view API:
// https://github.com/jenkinsci/pipeline-stage-view-plugin/tree/master/rest-api

package gojenkins

import (
	"fmt"
	"regexp"
)

var baseURLRegex *regexp.Regexp

func init() {
	var err error
	baseURLRegex, err = regexp.Compile("(.+)/wfapi/.*$")
	if err != nil {
		panic(err)
	}
}

type PipelineRun struct {
	Job       *Job
	Base      string
	URLs      map[string]map[string]string `json:"_links"`
	ID        string
	Name      string
	Status    string
	StartTime int64 `json:"startTimeMillis"`
	EndTime   int64 `json:"endTimeMillis"`
	Duration  int64 `json:"durationMillis"`
	Stages    []PipelineNode
}

type PipelineNode struct {
	Run            *PipelineRun
	Base           string
	URLs           map[string]map[string]string `json:"_links"`
	ID             string
	Name           string
	Status         string
	StartTime      int64 `json:"startTimeMillis"`
	Duration       int64 `json:"durationMillis"`
	StageFlowNodes []PipelineNode
	ParentNodes    []int64
}

type PipelineInputAction struct {
	ID         string
	Message    string
	ProceedURL string
	AbortURL   string
}

type PipelineArtifact struct {
	ID   string
	Name string
	Path string
	URL  string
	size int
}

type PipelineNodeLog struct {
	NodeID     string
	NodeStatus string
	Length     int64
	HasMore    bool
	Text       string
	ConsoleURL string
}

// utility function to fill in the Base fields under PipelineRun
func (run *PipelineRun) update() {
	href := run.URLs["self"]["href"]
	if matches := baseURLRegex.FindStringSubmatch(href); len(matches) > 1 {
		run.Base = matches[1]
	}
	for i := range run.Stages {
		run.Stages[i].Run = run
		href := run.Stages[i].URLs["self"]["href"]
		if matches := baseURLRegex.FindStringSubmatch(href); len(matches) > 1 {
			run.Stages[i].Base = matches[1]
		}
	}
}

func (job *Job) GetPipelineRuns() (pr []PipelineRun, err error) {
	_, err = job.Jenkins.Requester.GetJSON(job.Base+"/wfapi/runs", &pr, nil)
	if err != nil {
		return nil, err
	}
	for i := range pr {
		pr[i].update()
		pr[i].Job = job
	}

	return pr, nil
}

func (job *Job) GetPipelineRun(id string) (pr *PipelineRun, err error) {
	pr = new(PipelineRun)
	href := job.Base + "/" + id + "/wfapi/describe"
	_, err = job.Jenkins.Requester.GetJSON(href, pr, nil)
	if err != nil {
		return nil, err
	}
	pr.update()
	pr.Job = job

	return pr, nil
}

func (pr *PipelineRun) GetPendingInputActions() (PIAs []PipelineInputAction, err error) {
	PIAs = make([]PipelineInputAction, 0, 1)
	href := pr.Base + "/wfapi/pendingInputActions"
	_, err = pr.Job.Jenkins.Requester.GetJSON(href, &PIAs, nil)
	if err != nil {
		return nil, err
	}

	return PIAs, nil
}

func (pr *PipelineRun) GetArtifacts() (artifacts []PipelineArtifact, err error) {
	artifacts = make([]PipelineArtifact, 0, 0)
	href := pr.Base + "/wfapi/artifacts"
	_, err = pr.Job.Jenkins.Requester.GetJSON(href, artifacts, nil)
	if err != nil {
		return nil, err
	}

	return artifacts, nil
}

func (pr *PipelineRun) GetNode(id string) (node *PipelineNode, err error) {
	node = new(PipelineNode)
	href := pr.Base + "/execution/node/" + id + "/wfapi/describe"
	_, err = pr.Job.Jenkins.Requester.GetJSON(href, node, nil)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (node *PipelineNode) GetLog() (log *PipelineNodeLog, err error) {
	log = new(PipelineNodeLog)
	href := node.Base + "/wfapi/log"
	fmt.Println(href)
	_, err = node.Run.Job.Jenkins.Requester.GetJSON(href, log, nil)
	if err != nil {
		return nil, err
	}

	return log, nil
}
