package parsehub_go

import (
	"net/url"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"parsehub-go/internal"
	"time"
)

// ParseHub Run Wrapper
type Run struct {
	parsehub *ParseHub

	token    string
	response *RunResponse
	handler  RunHandler
}

// Get run data
func (r *Run) GetResponse() *RunResponse {
	return r.response
}

// This returns the data that was extracted by a run.
func (r *Run) LoadData(target interface{}) error {
	internal.Logf("Run: Load data for run %v", r.token)

	requestUrl, _ := url.Parse(ParseHubBaseUrl + "v2/runs/" + r.token + "/data")

	values := url.Values{}
	values.Add("api_key", r.parsehub.apiKey)

	requestUrl.RawQuery = values.Encode()

	if resp, err := http.Get(requestUrl.String()); err != nil {
		return err
	} else {
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		internal.Logf("Run: Load data body string: %s", body)
		json.Unmarshal(body, target)

		internal.Logf("Run: Load data unmarshaled: %v", target)

		return nil
	}
}

// This cancels a run and changes its status to cancelled. 
// Any data that was extracted so far will be available.
func (r *Run) Cancel() *Run {
	internal.Logf("Run: Cancel run %v", r.token)
	requestUrl, _ := url.Parse(ParseHubBaseUrl + "v2/runs/" + r.token + "/cancel")

	values := url.Values{}
	values.Add("api_key", r.parsehub.apiKey)
	requestUrl.RawQuery = values.Encode()

	if resp, err := http.PostForm(requestUrl.String(), values); err != nil {
		panic(err)
	} else {
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		internal.Logf("Run: Cancel run response string: %s", body)

		runResponse := &RunResponse{}
		json.Unmarshal(body, runResponse)

		internal.Logf("Run: Cancel run response: %+v", runResponse)

		r.response = runResponse // update response

		return r
	}
}

// This cancels a run if running, and deletes the run and its data.
func (r *Run) Delete() {
	internal.Logf("Run: Delete run %v", r.token)
	requestUrl, _ := url.Parse(ParseHubBaseUrl + "v2/runs/" + r.token)

	values := url.Values{}
	values.Add("api_key", r.parsehub.apiKey)
	requestUrl.RawQuery = values.Encode()

	request, _ := http.NewRequest(http.MethodDelete, requestUrl.String(), nil)

	if resp, err := http.DefaultClient.Do(request); err != nil {
		panic(err)
	} else {
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		internal.Logf("Run: Delete run response string: %s", body)

		runResponse := &RunResponse{}

		json.Unmarshal(body, runResponse)
		internal.Logf("Run: Delete run response: %v", runResponse)

		r.response = runResponse
	}
}

// Watch run
func (r *Run) Watch() {
	internal.Logf("Run: Start watching run %s", r.token)
	for {
		time.Sleep(10 * time.Second) // todo: delete hardcoded time

		internal.Logf("Run: Watch iteration for run %s", r.token)
		r.parsehub.GetRun(r.token)

		if r.response.DataReady > 0 {
			internal.Logf("Run: Watch closed. Handle run %s", r.token)
			r.handler.Handle(r)
			return
		}
	}
}
