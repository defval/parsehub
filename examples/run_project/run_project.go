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
	parsehub_go.SetLogger(logger)
	parsehub := parsehub_go.NewParseHub(parsehub_go.ApiKey)

	project := parsehub.GetProject(parsehub_go.ProjectToken)

	run := project.Run(parsehub_go.ProjectRunParams{
		StartTemplate: parsehub_go.StartTemplate,
		StartUrl: parsehub_go.StartUrl,
	}, &TestRunHandler{})

	fmt.Println(run)

	time.Sleep(50 * time.Second)
}

type TestRunHandler struct {

}

func (h *TestRunHandler) Handle(run *parsehub_go.Run) error {
	val := struct {
		Matches []struct {
			Id      string
			hltvUrl string
		}
	}{}
	run.LoadData(val)

	fmt.Println("result", val)
	return nil
}
