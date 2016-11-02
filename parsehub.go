package parsehub

import (
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"gopkg.in/mb24dev/parsehub.v1/internal"
)

const (
	ParseHubBaseUrl = "https://www.parsehub.com/api/"
)

// ParseHub adapter
type ParseHub struct {
	apiKey          string
	watchQueue      chan *Run
	projectRegistry map[string]*Project
	runRegistry     map[string]*Run
}

// Creates new ParseHub adapter with api key
func NewParseHub(apiKey string) *ParseHub {
	debugf("ParseHub: Create new parsehub client with api key: %v", apiKey)
	parsehub := &ParseHub{
		apiKey: apiKey,
		projectRegistry: map[string]*Project{},
		runRegistry: map[string]*Run{},
	}

	return parsehub
}

// This will return all of the projects in your account
func (parsehub *ParseHub) GetAllProjects() ([]*Project, error) {
	requestUrl, _ := url.Parse(ParseHubBaseUrl + "v2/projects")

	values := url.Values{}
	values.Add("api_key", parsehub.apiKey)

	requestUrl.RawQuery = values.Encode()

	if resp, err := http.Get(requestUrl.String()); err != nil {
		warningf("ParseHub.GetAllProjects: ParseHub HTTP request problem: %s", err.Error())
		return nil, err
	} else if success, err := internal.CheckHTTPStatusCode(resp.StatusCode); success {
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		debugf("ParseHub.GetAllProjects: Response string: %s", body)

		projectsResponse := &ProjectsResponse{}
		if err := json.Unmarshal(body, projectsResponse); err != nil {
			warningf("ParseHub.GetAllProjects: Unmarshal error with body %s", body)
			return nil, err
		}

		projects := []*Project{}
		var p *Project

		for _, projectResponse := range projectsResponse.Projects {
			p = NewProject(parsehub, projectResponse.Token)
			p.response = projectResponse
			projects = append(projects, p)
		}

		debugf("ParseHub.GetAllProjects: Get all projects response: %v", projects)

		return projects, nil
	} else {
		warningf("ParseHub.GetAllProjects: ParseHub HTTP response problem: %s", err.Error())
		return nil, err
	}

}

// This will return the project object wrapper for a specific project.
// 
// Params:
//
// start_url (Optional)	
// The url to start running on. Defaults to the project’s start_site.
//
// start_template (Optional)	
// The template to start running with.
// Defaults to the projects’s start_template (inside the options_json).
//
// start_value_override (Optional)	
// The starting global scope for this run. This can be used to pass parameters to your run. 
// For example, you can pass {"query": "San Francisco"} to use the query somewhere in your run. 
// Defaults to the project’s start_value.
//
// send_email (Optional)	
// If set to anything other than 0, send an email when the run either completes successfully 
// or fails due to an error. Defaults to 0.
func (parsehub *ParseHub) GetProject(projectToken string) (*Project, error) {
	debugf("ParseHub.GetProject: Get project with token: %s", projectToken)
	requestUrl, _ := url.Parse(ParseHubBaseUrl + "v2/projects/" + projectToken)

	values := url.Values{}
	values.Add("api_key", parsehub.apiKey)

	requestUrl.RawQuery = values.Encode()

	debugf("ParseHub.GetProject: Requested Url: %s", requestUrl)
	if resp, err := http.Get(requestUrl.String()); err != nil {
		warningf("ParseHub.GetProject: ParseHub HTTP request problem: %s", err.Error())
		return nil, err
	} else if success, err := internal.CheckHTTPStatusCode(resp.StatusCode); success {
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		projectResponse := &ProjectResponse{}

		if err := json.Unmarshal(body, projectResponse); err != nil {
			warningf("ParseHub.GetProject: Unmarshal error with body %s", body)
			return nil, err
		}

		internal.Lock.RLock()
		project := parsehub.projectRegistry[projectToken]
		internal.Lock.RUnlock()

		debugf("ParseHub.GetProject: Loaded project with token %s from registry: %+v", projectToken, project)

		if project == nil {
			debugf("ParseHub.GetProject: Need to put new project with token %s into registry", projectToken)
			project = NewProject(parsehub, projectToken)

			internal.Lock.RLock()
			parsehub.projectRegistry[projectToken] = project
			internal.Lock.RUnlock()
		}

		project.response = projectResponse

		return project, nil
	} else {
		warningf("ParseHub.GetProject: ParseHub HTTP response problem: %s", err.Error())
		return nil, err
	}
}

// This returns the run object wrapper for a given run token.
func (parsehub *ParseHub) GetRun(runToken string) (*Run, error) {
	debugf("ParseHub.GetRun: Get run with token %s", runToken)
	requestUrl, _ := url.Parse(ParseHubBaseUrl + "v2/runs/" + runToken)

	values := url.Values{}
	values.Add("api_key", parsehub.apiKey)

	requestUrl.RawQuery = values.Encode()

	if resp, err := http.Get(requestUrl.String()); err != nil {
		warningf("ParseHub.GetRun: ParseHub HTTP request problem: %s", err.Error())
		return nil, err
	} else if success, err := internal.CheckHTTPStatusCode(resp.StatusCode); success {
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		debugf("ParseHub.GetRun: Response string for run %s: %s", runToken, body)

		runResponse := &RunResponse{}
		if err := json.Unmarshal(body, runResponse); err != nil {
			warningf("ParseHub.GetRun: Problem with unmarshal json string: %s", body)
			return nil, err
		}

		debugf("ParseHub.GetRun: Run response: %+v", runResponse)

		internal.Lock.RLock()
		run := parsehub.runRegistry[runToken]
		internal.Lock.RUnlock()

		if run == nil {
			run = NewRun(parsehub, runToken)
			internal.Lock.Lock()
			parsehub.runRegistry[runToken] = run
			internal.Lock.Unlock()
		}

		run.response = runResponse // update response

		return run, nil
	} else {
		warningf("ParseHub.GetRun: ParseHub HTTP response problem: %s", err.Error())
		return nil, err
	}
}

// Loads run from string
// Example from webhook post body
func (parsehub *ParseHub) LoadRunFromBytes(body []byte) (*Run, error) {
	runResponse := &RunResponse{}

	if err := json.Unmarshal(body, runResponse); err != nil {
		warningf("NewRunFromString: Problem with unmarshal json string: %s", body)
		return nil, err
	}

	run := NewRun(parsehub, runResponse.RunToken)
	return run, nil
}
