package zap

import (
	"testing"

	"github.com/Jeffail/gabs/v2"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/afero"
)

const duplicatedYAML = `
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
zz:
  expand: secondaryexpansion.com
`

const badkeysYAML = `
e:
  bad_key: example.com
  a:
    expand: apples
g:
  expand: github.com
  d:
    expand: issmirnov/dotfiles
`

const badValuesYAML = `
e:
  expand: 2
  s:
    query: 3
g:
  expand: github.com
  ssl_off: "not_bool"
l:
  port: "not_int"
`

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
    me:
      expand: issmirnov
      z:
        expand: zap
    ak:
      query: apache/kafka
      c:
        query: +connect
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
ch:
  # expand: "/"
  v:
    expand: version # should expand to chrome://version
  'n':
    expand: net-internals
    d:
      expand: '#dns'
  schema: chrome 
`

func loadTestYaml() (*gabs.Container, error) {
	return parseYamlString(cYaml)
}

func TestParseYaml(t *testing.T) {
	Convey("Given a valid 'c.yml' file", t, func() {
		Afero = &afero.Afero{Fs: afero.NewMemMapFs()}
		Afero.WriteFile("c.yml", []byte(cYaml), 0644)
		c, err := ParseYaml("c.yml")
		Convey("ParseYaml should throw no error", func() {
			So(err, ShouldBeNil)
		})
		Convey("the gabs object should have path 'zz' present", func() {
			So(c.ExistsP("zz"), ShouldBeTrue)
		})
	})
}

func TestValidateConfig(t *testing.T) {
	Convey("Given a correctly formatted yaml Config", t, func() {
		conf, _ := parseYamlString(cYaml)
		//fmt.Printf(err.Error())
		Convey("The validator should pass", func() {
			So(ValidateConfig(conf), ShouldBeNil)
		})
	})

	// The YAML libraries don't have support for detecting duplicate keys
	// at parse time. Users will have to figure this out themselves.
	//Convey("Given a yaml Config with duplicated keys", t, func() {
	//	conf, _ := parseYamlString(duplicatedYAML)
	//	Convey("The validator should complain", func() {
	//		So(ValidateConfig(conf), ShouldNotBeNil)
	//	})
	//})

	Convey("Given a YAML Config with unknown keys", t, func() {
		conf, _ := parseYamlString(badkeysYAML)
		Convey("The validator should raise an error", func() {
			So(ValidateConfig(conf), ShouldNotBeNil)
		})
	})

	Convey("Given a YAML Config with malformed values", t, func() {
		conf, _ := parseYamlString(badValuesYAML)
		err := ValidateConfig(conf)
		Convey("The validator should raise a ton of errors", func() {
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "expected float64 value for string, got: not_int")
			So(err.Error(), ShouldContainSubstring, "expected string value for string, got: 3")
			So(err.Error(), ShouldContainSubstring, "expected bool value for string, got: not_bool")
			So(err.Error(), ShouldContainSubstring, "expected string value for string, got: 2")
		})
	})
}
