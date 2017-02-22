package main

import (
  "fmt"

  "github.com/Jeffail/gabs"
  "github.com/ghodss/yaml"
)

const cYaml = `
sp:
  expand: smirnov.wiki/project
  h:
    expand: hydra
sg:
  expand: smirnov.wiki/goals
  'n':
    expand: 2017
sd:
  expand: smirnov.wiki/device
  p:
    expand: puma
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
