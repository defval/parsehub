# ParseHub client [![Build Status](https://api.travis-ci.org/mb24dev/parsehub.svg?branch=master)](https://travis-ci.org/mb24dev/parsehub)

API docs: https://godoc.org/github.com/defval/parsehub.

Examples: https://godoc.org/github.com/defval/parsehub#pkg-examples.

## Installation

Install:

```shell
go get github.com/defval/parsehub
```

Import:

```go
import "github.com/defval/parsehub"
```

## Quickstart

```go
func ExampleParseHub_GetProjectAndRun() {
	parsehub := NewParseHub(ApiKey)

	if project, err := parsehub.GetProject(ProjectToken); err != nil {
		// handle error
	} else {
		// async run
		project.Run(ProjectRunParams{
			StartTemplate: StartTemplate,
			StartUrl: StartUrl,
		}, func(run *Run) error {
		
		    // handle run data
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
		})
	}

	// code that save main thread
}
```

