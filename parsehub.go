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

type ParseHub struct {
	apiKey          string
	watchQueue      chan *Run
	projectRegistry map[string]*Project
	runRegistry     map[string]*Run
}

func NewParseHub(apiKey string) *ParseHub {
	internal.Logf("ParseHub: Created new parsehub client with api key: %v", apiKey)
	parsehub := &ParseHub{
		apiKey: apiKey,
		projectRegistry: map[string]*Project{},
		runRegistry: map[string]*Run{},
	}

	return parsehub
}

// Set Logger
func SetLogger(logger *log.Logger) {
	internal.Logger = logger
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
	projects := &ProjectsResponse{}
	json.Unmarshal(body, projects)

	internal.Logf("ParseHub: Get all projects body: %v", string(body))
	internal.Logf("ParseHub: Get all projects response: %v", projects)

	return projects.Projects
}

// This will return the project object for a specific project.
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
	internal.Logf("ParseHub: Get project %s", projectToken)
	requestUrl, _ := url.Parse(ParseHubBaseUrl + "v2/projects/" + projectToken)

	values := url.Values{}
	values.Add("api_key", parsehub.apiKey)

	requestUrl.RawQuery = values.Encode()

	resp, _ := http.Get(requestUrl.String())

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	projectResponse := &ProjectResponse{}
	json.Unmarshal(body, projectResponse)

	project := parsehub.projectRegistry[projectToken]

	internal.Logf("ParseHub: Loaded project from registry: %+v", project)

	if project == nil {
		internal.Logf("ParseHub: Need to put project %s into registry", projectToken)
		project = &Project{}
		project.token = projectToken
		project.parsehub = parsehub
		parsehub.projectRegistry[projectToken] = project
	}

	project.response = projectResponse

	return project
}

func (parsehub *ParseHub) GetRun(runToken string) *Run {
	internal.Logf("ParseHub: Get run %s", runToken)
	requestUrl, _ := url.Parse(ParseHubBaseUrl + "v2/runs/" + runToken)

	values := url.Values{}
	values.Add("api_key", parsehub.apiKey)

	requestUrl.RawQuery = values.Encode()

	resp, _ := http.Get(requestUrl.String())

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	internal.Logf("ParseHub: Body data for run %s: %s", runToken, body)

	runResponse := &RunResponse{}
	err := json.Unmarshal(body, runResponse)

	if err != nil {
		panic(err)
	}

	internal.Logf("ParseHub: Run response: %+v", runResponse)

	run := parsehub.runRegistry[runToken]

	if run == nil {
		run := &Run{}
		run.token = runToken
		run.parsehub = parsehub
		parsehub.runRegistry[runToken] = run
	}

	run.token = runToken
	run.response = runResponse

	internal.Logf("ParseHub: Get run data: %+v", run.response)

	return run
}
