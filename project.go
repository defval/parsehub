package parsehub_go

import (
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
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
	data     *ProjectResponse
}

// Creates new project
func NewProject(projectToken string) *Project {
	return &Project{
		token: projectToken,
	}
}

// Get project data
func (p *Project) GetData() *ProjectResponse {
	return p.data
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
func (p *Project) Run(params ProjectRunParams) *Run {
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
			panic(err)
		} else {
			values.Add("start_value_override", string(bytes))
		}

	}

	if params.SendEmail {
		values.Add("send_email", "1")
	}

	if resp, err := http.PostForm(requestUrl.String(), values); err != nil {
		panic(err)
	} else {
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		runResponse := &RunResponse{}
		json.Unmarshal(body, runResponse)

		run := &Run{}
		run.parsehub = p.parsehub
		run.data = runResponse

		return run
	}
}

