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
	parsehub_go.SetLogger(logger)

	parsehub := parsehub_go.NewParseHub(parsehub_go.ApiKey)

	project := parsehub.GetProject(parsehub_go.ProjectToken)

	fmt.Printf("%+v", project)
}
