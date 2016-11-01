package main

import (
	"parsehub-go"
	"fmt"
	"log"
	"os"
)

func main() {
	logger := &log.Logger{}
	logger.SetOutput(os.Stdout)

	parsehub_go.SetLogger(parsehub_go.LogLevelDebug, logger)

	parsehub := parsehub_go.NewParseHub(parsehub_go.ApiKey)

	if project, err := parsehub.GetProject(parsehub_go.ProjectToken); err != nil {
		log.Fatalf(err.Error())
	} else {
		fmt.Printf("%+v", project)
	}
}
