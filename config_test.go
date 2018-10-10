package main

import (
	"testing"

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

func TestParseYaml(t *testing.T) {
	Convey("Given a valid 'c.yml' file", t, func() {
		Afero = &afero.Afero{Fs: afero.NewMemMapFs()}
		Afero.WriteFile("c.yml", []byte(cYaml), 0644)
		c, err := parseYaml("c.yml")
		Convey("parseYaml should throw no error", func() {
			So(err, ShouldBeNil)
		})
		Convey("the gabs object should have path 'zz' present", func() {
			So(c.ExistsP("zz"), ShouldBeTrue)
		})
	})
}

func TestValidateConfig(t *testing.T) {
	Convey("Given a correctly formatted yaml config", t, func() {
		conf, _ := parseYamlString(cYaml)
		//fmt.Printf(err.Error())
		Convey("The validator should pass", func() {
			So(validateConfig(conf), ShouldBeNil)
		})
	})

	// The YAML libraries don't have support for detecting duplicate keys
	// at parse time. Users will have to figure this out themselves.
	//Convey("Given a yaml config with duplicated keys", t, func() {
	//	conf, _ := parseYamlString(duplicatedYAML)
	//	Convey("The validator should complain", func() {
	//		So(validateConfig(conf), ShouldNotBeNil)
	//	})
	//})

	Convey("Given a YAML config with unknown keys", t, func() {
		conf, _ := parseYamlString(badkeysYAML)
		Convey("The validator should raise an error", func() {
			So(validateConfig(conf), ShouldNotBeNil)
		})
	})
}
