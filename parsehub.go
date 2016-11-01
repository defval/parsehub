package parsehub_go

import (
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"parsehub-go/internal"
	"log"
)

const (
	ParseHubBaseUrl = "https://www.parsehub.com/api/"
)

// Set Logger for package
func SetLogger(logger *log.Logger) {
	internal.Logger = logger
}

// ParseHub adapter
type ParseHub struct {
	apiKey          string
	watchQueue      chan *Run
	projectRegistry map[string]*Project
	runRegistry     map[string]*Run
}

// Creates new ParseHub adapter with api key
func NewParseHub(apiKey string) *ParseHub {
	internal.Logf("ParseHub: Create new parsehub client with api key: %v", apiKey)
	parsehub := &ParseHub{
		apiKey: apiKey,
		projectRegistry: map[string]*Project{},
		runRegistry: map[string]*Run{},
	}

	return parsehub
}

// This will return all of the projects in your account
func (parsehub *ParseHub) GetAllProjects() []*ProjectResponse {
	requestUrl, _ := url.Parse(ParseHubBaseUrl + "v2/projects")

	values := url.Values{}
	values.Add("api_key", parsehub.apiKey)

	requestUrl.RawQuery = values.Encode()

	resp, _ := http.Get(requestUrl.String())

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	internal.Logf("ParseHub.GetAllProjects: Response string: %s", body)

	projects := &ProjectsResponse{}
	json.Unmarshal(body, projects)
	internal.Logf("ParseHub.GetAllProjects: Get all projects response: %v", projects)

	return projects.Projects
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
func (parsehub *ParseHub) GetProject(projectToken string) *Project {
	internal.Logf("ParseHub.GetProject: Get project with token: %s", projectToken)
	requestUrl, _ := url.Parse(ParseHubBaseUrl + "v2/projects/" + projectToken)

	values := url.Values{}
	values.Add("api_key", parsehub.apiKey)

	requestUrl.RawQuery = values.Encode()

	if resp, err := http.Get(requestUrl.String()); err != nil {
		panic(err)
	} else {
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		projectResponse := &ProjectResponse{}
		json.Unmarshal(body, projectResponse)

		internal.Lock.RLock()
		project := parsehub.projectRegistry[projectToken]
		internal.Lock.RUnlock()

		internal.Logf("ParseHub.GetProject: Loaded project with token %s from registry: %+v", projectToken, project)

		if project == nil {
			internal.Logf("ParseHub.GetProject: Need to put new project with token %s into registry", projectToken)
			project = NewProject(parsehub, projectToken)

			internal.Lock.RLock()
			parsehub.projectRegistry[projectToken] = project
			internal.Lock.RUnlock()
		}

		project.response = projectResponse

		return project
	}
}

// This returns the run object wrapper for a given run token.
func (parsehub *ParseHub) GetRun(runToken string) *Run {
	internal.Logf("ParseHub.GetRun: Get run with token %s", runToken)
	requestUrl, _ := url.Parse(ParseHubBaseUrl + "v2/runs/" + runToken)

	values := url.Values{}
	values.Add("api_key", parsehub.apiKey)

	requestUrl.RawQuery = values.Encode()

	if resp, err := http.Get(requestUrl.String()); err != nil {
		panic(err) // todo: remove panic
	} else {
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		internal.Logf("ParseHub.GetRun: Response string for run %s: %s", runToken, body)

		runResponse := &RunResponse{}
		err := json.Unmarshal(body, runResponse)

		if err != nil {
			internal.Logf("ParseHub.GetRun: Problem with unmarshal json string: %s", body)
			panic(err) // todo: remove panic
		}

		internal.Logf("ParseHub.GetRun: Run response: %+v", runResponse)

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

		return run
	}
}
