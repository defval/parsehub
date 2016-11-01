package parsehub

// ParseHub Projects
type ProjectsResponse struct {
	Projects []*ProjectResponse `json:"projects"`
}

// ParseHub Project
type ProjectResponse struct {
	// A globally unique id representing this project.
	Token         string `json:"token"`

	// The title give by the user when creating this project.
	Title         string `json:"title"`

	// The JSON-stringified representation of all the instructions for running this project. 
	// This representation is not yet documented, but will eventually allow developers to create 
	// plugins for ParseHub.
	TemplatesJSON string `json:"templates_json"`

	// The name of the template with which ParseHub should start executing the project.
	MainTemplate  string `json:"main_template"`

	// The default URL at which ParseHub should start running the project.
	Main_site     string `json:"main_site"`

	// An object containing several advanced options for the project.
	OptionsJSON   string `json:"option_json"`

	// The run object of the most recently started run (orderd by start_time) for the project.
	LastRun       *RunResponse `json:"last_run"`

	// The run object of the most recent ready run (ordered by start_time) for the project. A ready run is one 
	// whose data_ready attribute is truthy. The last_run and last_ready_run for a project may be the same.
	LastReadyRun  *RunResponse `json:"last_ready_run"`
}


// ParseHub Run
type RunResponse struct {
	// A globally unique id representing the project that this run belongs to.
	ProjectToken  string `json:"project_token"`

	// A globally unique id representing this run.
	RunToken      string `json:"run_token"`

	// The status of the run. It can be one of initialized, queued, running, cancelled, complete, or error.
	Status        string `json:"status"`

	// Whether the data for this run is ready to download. If the status is complete, this will always be truthy. 
	// If the status is cancelled or error, then this may be truthy or falsy, depending on whether any 
	// data is available.
	DataReady     uint8 `json:"data_ready"`

	// The time that this run was started at, in UTC +0000.
	StartTime     string `json:"start_time"`

	// The time that this run was stopped. This field will be null if the run is either initialized or running. 
	// Time is in UTC +0000.
	EndTime       string `json:"end_time"`

	// The number of pages that have been traversed by this run so far.
	Pages         interface{} `json:"pages"` // todo: fix parsehub format

	// The md5sum of the results. This can be used to check if any results data has changed between two runs.
	Md5sum        string `json:"md5sum"`

	// The url that this run was started on.
	StartURL      string `json:"start_url"`

	// The template that this run was started with.
	StartTemplate string `json:"start_template"`

	// The starting value of the global scope for this run.
	StartValue    string `json:"start_value"`
}

