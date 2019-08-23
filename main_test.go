package main

import (
	"fmt"

	"github.com/Jeffail/gabs/v2"
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
z:
  expand: zero.com
  ssl_off: yes
zz:
  expand: zero.ssl.on.com
  ssl_off: no
l:
  expand: localhost
  ssl_off: yes
  a:
    port: 8080
    s:
      expand: service
ak:
  expand: kafka.apache.org
  hi:
    expand: contact
  "*":
    d:
      expand: documentation.html
    j:
      expand: javadoc/index.html?overview-summary.html
wc:
  expand: wildcard.com
  "*":
    "*":
      "*":
        four:
          expand: "4"
`

func loadTestYaml() (*gabs.Container, error) {
	return parseYamlString(cYaml)
}

func parseYamlString(config string) (*gabs.Container, error) {
	d, jsonErr := yaml.YAMLToJSON([]byte(config))
	if jsonErr != nil {
		fmt.Printf("Error encoding input to JSON.\n%s\n", jsonErr.Error())
		return nil, jsonErr
	}
	j, _ := gabs.ParseJSON(d)
	return j, nil
}
