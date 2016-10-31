package main

import (
	"parsehub-go"
	"fmt"
)

func main() {
	parsehub := parsehub_go.NewParseHub(parsehub_go.ApiKey)

	projects := parsehub.GetAllProjects()

	fmt.Printf("%+v", projects[0])
}
