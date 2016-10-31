package main

import (
	"parsehub-go"
	"fmt"
)

func main() {
	parsehub := parsehub_go.NewParseHub(parsehub_go.ApiKey)

	project := parsehub.GetProject(parsehub_go.ProjectToken)

	fmt.Println(project.GetData())

	run := project.Run(parsehub_go.ProjectRunParams{
		StartTemplate: parsehub_go.StartTemplate,
		StartUrl: parsehub_go.StartUrl,
	})

	fmt.Println(run.GetData())
}
