package parsehub_go

import (
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

const (
	ParseHubBaseUrl = "https://www.parsehub.com/api/"
)

type ParseHub struct {
	apiKey string
}

func NewParseHub(apiKey string) *ParseHub {
	return &ParseHub{
		apiKey: apiKey,
	}
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
	requestUrl, _ := url.Parse(ParseHubBaseUrl + "v2/projects/" + projectToken)

	values := url.Values{}
	values.Add("api_key", parsehub.apiKey)

	requestUrl.RawQuery = values.Encode()

	resp, _ := http.Get(requestUrl.String())

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	projectResponse := &ProjectResponse{}
	json.Unmarshal(body, projectResponse)

	project := NewProject(projectToken)
	project.parsehub = parsehub

	project.data = projectResponse

	return project
}
