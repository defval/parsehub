package parsehub

import (
	"fmt"
	"log"
	"os"
)

// Set parsehub library logger
func ExampleSetLogger() {
	logger := &log.Logger{}
	logger.SetOutput(os.Stdout)

	SetLogger(LogLevelDebug, logger)
}

// Run parsehub project with params and handle data async
// with polling status
func ExampleProject_Run() {
	parsehub := NewParseHub("__API_KEY__")

	handleFunc := func(run *Run) error {
		val := map[string]interface{}{}

		if err := run.LoadData(&val); err != nil {
			log.Fatalf(err.Error())
		}

		fmt.Println("result", val)

		// delete after extract data
		if err := run.Delete(); err != nil {
			log.Fatalf(err.Error())
		}
		return nil
	}

	if project, err := parsehub.GetProject("__PROJECT_TOKEN__"); err != nil {
		// handle error
	} else {
		// async run
		project.Run(ProjectRunParams{
			StartTemplate: "__START_TEMPLATE__",
			StartUrl:      "__START_URL__",
		}, handleFunc) // or nil if you use webhooks
	}

	// your code
}

// Load data from string
func ExampleParseHub_LoadRunFromBytes() {
	runJsonBytes := []byte("__RUN_FROM_WEBHOOK__")

	parsehub := NewParseHub("__API_KEY__")

	if run, err := parsehub.LoadRunFromBytes(runJsonBytes); err != nil {
		log.Fatalf(err.Error())
	} else {
		val := map[string]interface{}{}

		if err := run.LoadData(&val); err != nil {
			log.Fatalf(err.Error())
		}

		fmt.Println("result", val)

		// delete after extract data
		if err := run.Delete(); err != nil {
			log.Fatalf(err.Error())
		}
	}

}

// Load run data
func ExampleRun_LoadData() {
	parsehub := NewParseHub("__API_KEY__")

	if run, err := parsehub.GetRun("__RUN_TOKEN__"); err != nil {
		// handle error
	} else {
		v := map[string]interface{}{} // example struct
		run.LoadData(v)               // load data
		fmt.Printf("%+v", v)
	}
}

// Get parsehub project
func ExampleParseHub_GetProject() {
	parsehub := NewParseHub("__API_KEY__")

	if project, err := parsehub.GetProject("__PROJECT_TOKEN__"); err != nil {
		// handle error
	} else {
		fmt.Printf("%+v", project)
	}
}

// Get all parsehub project
func ExampleParseHub_GetAllProjects() {
	parsehub := NewParseHub("__API_KEY__")

	if projects, err := parsehub.GetAllProjects(); err != nil {
		// handle error
	} else {
		for _, project := range projects {
			fmt.Printf("%+v", project)
		}
	}
}

// Get parsehub run
func ExampleParseHub_GetRun() {
	parsehub := NewParseHub("__API_KEY__")

	if run, err := parsehub.GetRun("__RUN_TOKEN__"); err != nil {
		// handle error
	} else {
		fmt.Printf("%+v", run)
	}
}
