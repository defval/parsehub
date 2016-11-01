# ParseHub client [![Build Status](https://api.travis-ci.org/mb24dev/parsehub.svg?branch=master)](https://travis-ci.org/mb24dev/parsehub)

API docs: https://godoc.org/gopkg.in/mb24dev/parsehub.v1.

Examples: https://godoc.org/gopkg.in/mb24dev/parsehub.v1#pkg-examples.

## Installation

Install:

```shell
go get gopkg.in/mb24dev/parsehub.v1
```

Import:

```go
import "gopkg.in/mb24dev/parsehub.v1"
```

## Vendoring

If you are using a vendoring tool with support for semantic versioning
e.g. [glide](https://github.com/Masterminds/glide), you can import this
package via its GitHub URL:

```yaml
- package: github.com/mb24dev/parsehub
  version: ^1.0.0
```

WARNING: please note that by importing `github.com/mb24dev/parsehub`
directly (without semantic versioning constrol) you are in danger of
running in the breaking API changes. Use carefully and at your own
risk!

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

