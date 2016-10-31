package main

import (
	"parsehub-go"
	"fmt"
)

func main() {
	parsehub := parsehub_go.NewParseHub(parsehub_go.ApiKey)

	project := parsehub.GetProject(parsehub_go.ProjectToken)

	fmt.Printf("%+v", project)
}
