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

	projects := parsehub.GetAllProjects()

	for _, project := range projects {
		fmt.Printf("%+v", project)
	}
}
