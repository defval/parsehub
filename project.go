package parsehub

import (
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"gopkg.in/mb24dev/parsehub.v1/internal"
)

// Project run params
type ProjectRunParams struct {
	StartUrl           string
	StartTemplate      string
	StartValueOverride map[string]interface{}
	SendEmail          bool
}

// ParseHub project Wrapper
type Project struct {
	parsehub *ParseHub

	token    string
	response *ProjectResponse
}

// Creates new parsehub project wrapper
func NewProject(parsehub *ParseHub, token string) *Project {
	return &Project{
		parsehub: parsehub,
		token: token,
	}
}

// Get project data
func (p *Project) GetResponse() *ProjectResponse {
	return p.response
}

// This will start running an instance of the project on the ParseHub cloud. It will create a new run object. 
// This method will return immediately, while the run continues in the background. 
// You can use webhooks or polling to figure out when the data for this 
// run is ready in order to retrieve it.
// 
// Params:
// start_url (Optional)	
// The url to start running on. Defaults to the project’s start_site.
//
// start_template (Optional)	
// The template to start running with. Defaults to the projects’s start_template (inside the options_json).
//
// start_value_override (Optional)	
// The starting global scope for this run. This can be used to pass parameters to your run. 
// For example, you can pass {"query": "San Francisco"} to use the query somewhere in 
// your run. Defaults to the project’s start_value.
//
// send_email (Optional)	
// If set to anything other than 0, send an email when the run either completes successfully or 
// fails due to an error. Defaults to 0.
func (p *Project) Run(params ProjectRunParams, handleFunc HandleRunFunc) (*Run, error) {
	debugf(
		"Project.Run: Run project %s with params: %+v",
		p.token,
		params,
	)

	requestUrl, _ := url.Parse(ParseHubBaseUrl + "v2/projects/" + p.token + "/run")

	values := url.Values{}
	values.Add("api_key", p.parsehub.apiKey)

	if params.StartUrl != "" {
		values.Add("start_url", params.StartUrl)
	}

	if params.StartTemplate != "" {
		values.Add("start_template", params.StartTemplate)
	}

	if len(params.StartValueOverride) != 0 {
		if bytes, err := json.Marshal(params.StartValueOverride); err != nil {
			fatalf("Project.Run: Incorrect StartValueOverride")
		} else {
			values.Add("start_value_override", string(bytes))
		}
	}

	if params.SendEmail {
		values.Add("send_email", "1")
	}

	if resp, err := http.PostForm(requestUrl.String(), values); err != nil {
		warningf("Project.Run: ParseHub HTTP request problem: %s", err.Error())
		return nil, err
	} else if success, err := internal.CheckHTTPStatusCode(resp.StatusCode); success {
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		runResponse := &RunResponse{}
		if err := json.Unmarshal(body, runResponse); err != nil {
			warningf("Project.Run: Unmarshal error with body %s", body)
			return nil, err
		}

		internal.Lock.RLock()
		run := p.parsehub.runRegistry[runResponse.RunToken]
		internal.Lock.RUnlock()

		if run == nil {
			run = NewRun(p.parsehub, runResponse.RunToken)
			internal.Lock.Lock()
			p.parsehub.runRegistry[runResponse.RunToken] = run
			internal.Lock.Unlock()
		}

		run.response = runResponse

		run.SetHandler(handleFunc)

		internal.Lock.Lock()
		p.parsehub.runRegistry[run.token] = run
		internal.Lock.Unlock()

		debugf("Project.Run: Start WatchAndHandle for run with token %s", run.token)
		go run.WatchAndHandle()

		return run, nil
	} else {
		warningf("Project.Run: ParseHub HTTP response problem: %s", err.Error())
		return nil, err
	}
}

// This returns the data for the most recent ready run for a project. 
// You can use this method in order to have a synchronous interface to your project.
func (p *Project) LoadLastReadyData(target interface{}) error {
	debugf("Project.LoadLastReadyData: Load: %s", p.token)
	requestUrl, _ := url.Parse(ParseHubBaseUrl + "v2/projects/" + p.token + "/last_ready_run/data")

	values := url.Values{}
	values.Add("api_key", p.parsehub.apiKey)

	requestUrl.RawQuery = values.Encode()

	if resp, err := http.Get(requestUrl.String()); err != nil {
		warningf("Project.LoadLastReadyData: ParseHub HTTP request problem: %s", err.Error())
		return err
	} else if success, err := internal.CheckHTTPStatusCode(resp.StatusCode); success {
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(body, target); err != nil {
			warningf("Project.LoadLastReadyData: Unmarshal error with body %s", body)
			return err
		}
		return nil
	} else {
		warningf("Project.LoadLastReadyData: ParseHub HTTP response problem: %s", err.Error())
		return err
	}
}

