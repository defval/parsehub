package parsehub

import (
	"net/url"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"time"
	"gopkg.in/mb24dev/parsehub.v1/internal"
)

type HandleRunFunc func(run *Run) error

// ParseHub Run Wrapper
type Run struct {
	parsehub   *ParseHub

	token      string
	response   *RunResponse
	handleFunc HandleRunFunc

	watching   bool
}

// Creates new ParseHub run wrapper
func NewRun(parsehub *ParseHub, token string) *Run {
	return &Run{
		parsehub: parsehub,
		token: token,
	}
}

// Set run handler
func (r *Run) SetHandler(handleFunc HandleRunFunc) {
	r.handleFunc = handleFunc
}

// Get run data
func (r *Run) GetResponse() *RunResponse {
	return r.response
}

// This load the data that was extracted by a run.
func (r *Run) LoadData(target interface{}) error {
	debugf("Run.LoadData: Load data for run %v", r.token)

	requestUrl, _ := url.Parse(ParseHubBaseUrl + "v2/runs/" + r.token + "/data")

	values := url.Values{}
	values.Add("api_key", r.parsehub.apiKey)

	requestUrl.RawQuery = values.Encode()

	if resp, err := http.Get(requestUrl.String()); err != nil {
		warningf("Run.LoadData: ParseHub HTTP request problem: %s", err.Error())
		return err
	} else if success, err := internal.CheckHTTPStatusCode(resp.StatusCode); success {
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		debugf("Run.LoadData: Body string: %s", body)

		if err := json.Unmarshal(body, target); err != nil {
			warningf("Run.LoadData: Unmarshal error with body %s", body)
			return err
		}

		debugf("Run.LoadData: Target: %+v", target)
		return err
	} else {
		warningf("Run.LoadData: ParseHub HTTP response problem: %s", err.Error())
		return err
	}
}

// This cancels a run and changes its status to cancelled. 
// Any data that was extracted so far will be available.
func (r *Run) Cancel() error {
	debugf("Run.Cancel: Cancel run %v", r.token)
	requestUrl, _ := url.Parse(ParseHubBaseUrl + "v2/runs/" + r.token + "/cancel")

	values := url.Values{}
	values.Add("api_key", r.parsehub.apiKey)
	requestUrl.RawQuery = values.Encode()

	if resp, err := http.PostForm(requestUrl.String(), values); err != nil {
		warningf("Run.Cancel: ParseHub HTTP request problem: %s", err.Error())
		return err
	} else if success, err := internal.CheckHTTPStatusCode(resp.StatusCode); success {
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		debugf("Run.Cancel: Cancel run response string: %s", body)

		runResponse := &RunResponse{}
		if err := json.Unmarshal(body, runResponse); err != nil {
			warningf("Run.Cancel: Unmarshal error with body %s", body)
			return err
		}

		debugf("Run.Cancel: Cancel run response: %+v", runResponse)

		r.response = runResponse // update response

		return nil
	} else {
		warningf("Run.Cancel: ParseHub HTTP response problem: %s", err.Error())
		return err
	}

}

// Refresh run data
func (r *Run) Refresh() error {
	_, err := r.parsehub.GetRun(r.token)
	return err
}

// This cancels a run if running, and deletes the run and its data.
func (r *Run) Delete() error {
	debugf("Run.Delete: Delete run %v", r.token)
	requestUrl, _ := url.Parse(ParseHubBaseUrl + "v2/runs/" + r.token)

	values := url.Values{}
	values.Add("api_key", r.parsehub.apiKey)
	requestUrl.RawQuery = values.Encode()

	request, _ := http.NewRequest(http.MethodDelete, requestUrl.String(), nil)

	if resp, err := http.DefaultClient.Do(request); err != nil {
		warningf("Run.Delete: ParseHub HTTP request problem: %s", err.Error())
		return err
	} else if success, err := internal.CheckHTTPStatusCode(resp.StatusCode); success {
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		debugf("Run.Delete: Delete run response string: %s", body)

		runResponse := &RunResponse{}

		if err := json.Unmarshal(body, runResponse); err != nil {
			warningf("Run.Delete: Unmarshal error with body %s", body)
			return err
		}
		debugf("Run.Delete: Delete run response: %v", runResponse)

		r.response = runResponse

		internal.Lock.Lock()
		delete(r.parsehub.runRegistry, r.token)
		internal.Lock.Unlock()
		return nil
	} else {
		warningf("Run.Delete: ParseHub HTTP response problem: %s", err.Error())
		return err
	}
}

// Watch for complete run and handle if handler exist
// Use SetHandler() for handle run data
func (r *Run) WatchAndHandle() {
	// No double watches
	if r.watching {
		warningf("Run.WatchAndHandle: Watching double run with token %s", r.token) // its not a problem
		return
	}

	debugf("Run.WatchAndHandle: Start watching run with token %s", r.token)
	r.watching = true

	for {
		time.Sleep(10 * time.Second) // todo: delete hardcoded time

		debugf("Run.WatchAndHandle: Watch iteration run with token %s", r.token)
		r.parsehub.GetRun(r.token)

		// todo: add conditions for stop watching
		if r.response.EndTime != "" {
			r.watching = false

			debugf("Run.WatchAndHandle: Watch finished. Handle run with token %s", r.token)
			if err := r.handleFunc(r); err != nil {
				warningf("Run.WatchAndHandle: Handle run with token %s error: %s", r.token, err.Error())
				return
			} else {
				internal.Lock.Lock()
				delete(r.parsehub.runRegistry, r.token)
				internal.Lock.Unlock()
			}

			return // stop watching
		}
	}
}
