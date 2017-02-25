package main

import (
	"fmt"

	"github.com/Jeffail/gabs"
	"github.com/ghodss/yaml"
)

const cYaml = `
e:
  expand: example.com
  a:
    expand: apples
  b:
    expand: bananas
g:
  expand: github.com
  d:
    expand: issmirnov/dotfiles
  z:
    expand: issmirnov/zap
  s:
    query: "search?q="

`

func parseDummyYaml() (*gabs.Container, error) {
	d, jsonErr := yaml.YAMLToJSON([]byte(cYaml))
	if jsonErr != nil {
		fmt.Printf("Error encoding input to JSON.\n%s\n", jsonErr.Error())
		return nil, jsonErr
	}
	j, _ := gabs.ParseJSON(d)
	return j, nil

}
