package main

import (
	"parsehub-go"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	logger := &log.Logger{}
	logger.SetOutput(os.Stdout)

	parsehub_go.SetLogger(parsehub_go.LogLevelDebug, logger)

	parsehub := parsehub_go.NewParseHub(parsehub_go.ApiKey)

	runHandler := &TestRunHandler{}

	if project, err := parsehub.GetProject(parsehub_go.ProjectToken); err != nil {
		log.Fatalf(err.Error())
	} else {
		// concurrent watches
		project.Run(parsehub_go.ProjectRunParams{
			StartTemplate: parsehub_go.StartTemplate,
			StartUrl: parsehub_go.StartUrl,
		}, runHandler.Handle)

		project.Run(parsehub_go.ProjectRunParams{
			StartTemplate: parsehub_go.StartTemplate,
			StartUrl: parsehub_go.StartUrl,
		}, runHandler.Handle)

		// other run
		project.Run(parsehub_go.ProjectRunParams{
			StartTemplate: parsehub_go.StartTemplate2,
			StartUrl: parsehub_go.StartUrl2,
		}, runHandler.Handle)

		project.Run(parsehub_go.ProjectRunParams{
			StartTemplate: parsehub_go.StartTemplate2,
			StartUrl: parsehub_go.StartUrl2,
		}, runHandler.Handle)

		project.Run(parsehub_go.ProjectRunParams{
			StartTemplate: parsehub_go.StartTemplate2,
			StartUrl: parsehub_go.StartUrl2,
		}, runHandler.Handle)
	}

	time.Sleep(100 * time.Second) // hard code for save main thread
}

// This struct handle run completed run
type TestRunHandler struct {

}

func (h *TestRunHandler) Handle(run *parsehub_go.Run) error {
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
