package main

import (
	"fmt"
	"io/ioutil"

	"github.com/Jeffail/gabs"
	"github.com/ghodss/yaml"
)

// expandKey is used to get expand mapping
const expandKey string = "expand"

// queryKey is used to get search prefix
const queryKey string = "query"

// parseYaml takes a file name and returns a gabs config object.
func parseYaml(fname string) (*gabs.Container, error) {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		fmt.Printf("Unable to read file: %s\n", err.Error())
		return nil, err
	}
	d, jsonErr := yaml.YAMLToJSON([]byte(data))
	if jsonErr != nil {
		fmt.Printf("Error encoding input to JSON.\n%s\n", jsonErr.Error())
		return nil, jsonErr
	}
	j, _ := gabs.ParseJSON(d)
	return j, nil
}
