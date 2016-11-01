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

	watching bool
}

// Creates new ParseHub run wrapper
func NewRun(parsehub *ParseHub, token string) *Run {
	return &Run{
		parsehub: parsehub,
		token: token,
	}
}

// Set run handler
func (r *Run) SetHandler(handler RunHandler) {
	r.handler = handler
}

// Get run data
func (r *Run) GetResponse() *RunResponse {
	return r.response
}

// This load the data that was extracted by a run.
func (r *Run) LoadData(target interface{}) error {
	internal.Logf("Run.LoadData: Load data for run %v", r.token)

	requestUrl, _ := url.Parse(ParseHubBaseUrl + "v2/runs/" + r.token + "/data")

	values := url.Values{}
	values.Add("api_key", r.parsehub.apiKey)

	requestUrl.RawQuery = values.Encode()

	if resp, err := http.Get(requestUrl.String()); err != nil {
		return err
	} else {
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		internal.Logf("Run.LoadData: Body string: %s", body)

		json.Unmarshal(body, target)
		internal.Logf("Run.LoadData: Target: %+v", target)

		return nil
	}
}

// This cancels a run and changes its status to cancelled. 
// Any data that was extracted so far will be available.
func (r *Run) Cancel() *Run {
	internal.Logf("Run.Cancel: Cancel run %v", r.token)
	requestUrl, _ := url.Parse(ParseHubBaseUrl + "v2/runs/" + r.token + "/cancel")

	values := url.Values{}
	values.Add("api_key", r.parsehub.apiKey)
	requestUrl.RawQuery = values.Encode()

	if resp, err := http.PostForm(requestUrl.String(), values); err != nil {
		panic(err)
	} else {
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		internal.Logf("Run.Cancel: Cancel run response string: %s", body)

		runResponse := &RunResponse{}
		json.Unmarshal(body, runResponse)

		internal.Logf("Run.Cancel: Cancel run response: %+v", runResponse)

		r.response = runResponse // update response

		return r
	}
}

// Refresh run data
func (r *Run) Refresh() {
	r.parsehub.GetRun(r.token)
}

// This cancels a run if running, and deletes the run and its data.
func (r *Run) Delete() {
	internal.Logf("Run.Delete: Delete run %v", r.token)
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
		internal.Logf("Run.Delete: Delete run response string: %s", body)

		runResponse := &RunResponse{}

		json.Unmarshal(body, runResponse)
		internal.Logf("Run.Delete: Delete run response: %v", runResponse)

		r.response = runResponse

		internal.Lock.Lock()
		delete(r.parsehub.runRegistry, r.token)
		internal.Lock.Unlock()
	}
}

// Watch for complete run and handle if handler exist
// Use SetHandler() for handle run data
func (r *Run) WatchAndHandle() {
	// No double watches
	if r.watching {
		internal.Logf("Run.WatchAndHandle: Watching double run with token %s", r.token)
		return
	}

	internal.Logf("Run.WatchAndHandle: Start watching run with token %s", r.token)
	r.watching = true

	for {
		time.Sleep(10 * time.Second) // todo: delete hardcoded time

		internal.Logf("Run.WatchAndHandle: Watch iteration run with token %s", r.token)
		r.parsehub.GetRun(r.token)

		// todo: add conditions for stop watching
		if r.response.EndTime != "" {
			r.watching = false

			internal.Logf("Run.WatchAndHandle: Watch finished. Handle run with token %s", r.token)
			if err := r.handler.Handle(r); err != nil {
				internal.Logf("Run.WatchAndHandle: Handle run with token %s error: %s", r.token, err.Error())
			} else {
				internal.Lock.Lock()
				delete(r.parsehub.runRegistry, r.token)
				internal.Lock.Unlock()
			}

			return // stop watching
		}
	}
}
