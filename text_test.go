package main

import (
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTokenizer(t *testing.T) {
	Convey("Given a string 'sp/h/monitoring'", t, func() {
		l := tokenize("sp/h/monitoring")
		Convey("The resulting list should have length 3", func() {
			So(l.Len(), ShouldEqual, 3)
		})
	})

	Convey("Given a string 'sp/h'", t, func() {
		l := tokenize("sp/h")
		Convey("The resulting list should", func() {

			Convey("Have length 2", func() {
				So(l.Len(), ShouldEqual, 2)
			})
			Convey("The first element should be equal to 'sp'", func() {
				So(l.Front().Value, ShouldEqual, "sp")
			})
			Convey("The last element should be equal to 'h'", func() {
				So(l.Back().Value, ShouldEqual, "h")
			})
		})
	})

	Convey("Given a string 'sg/n'", t, func() {
		l := tokenize("sg/n")
		Convey("The resulting list should", func() {

			Convey("Have length 2", func() {
				So(l.Len(), ShouldEqual, 2)
			})
			Convey("The first element should be equal to 'sg'", func() {
				So(l.Front().Value, ShouldEqual, "sg")
			})
			Convey("The last element should be equal to 'g'", func() {
				So(l.Back().Value, ShouldEqual, "n")
			})
		})
	})
}

func TestExpander(t *testing.T) {

	Convey("Given 'sp/h'", t, func() {
		c, _ := parseDummyYaml()
		l := tokenize("sp/h")
		var res bytes.Buffer
		res.WriteString("https:/")

		expand(c, l.Front(), &res)

		Convey("result should equal 'https://smirnov.wiki/project/hydra'", func() {
			So(res.String(), ShouldEqual, "https://smirnov.wiki/project/hydra")
		})
	})
	// Convey("Given 'sg/n'", t, func() {
	//     c, _ := parseDummyYaml()
	//     l := tokenize("sg/n ")
	//     var res bytes.Buffer
	//     res.WriteString("https:/")
	//
	//     expand(c, l.Front(), &res)
	//
	//     Convey("result should equal 'https://smirnov.wiki/goals/2017'", func() {
	//         So(res.String(), ShouldEqual, "https://smirnov.wiki/goals/2017")
	//     })
	// })
	Convey("Given 'sp/h/monitoring'", t, func() {
		c, _ := parseDummyYaml()
		l := tokenize("sp/h/monitoring")
		var res bytes.Buffer
		res.WriteString("https:/")

		expand(c, l.Front(), &res)

		Convey("result should equal 'https://smirnov.wiki/project/hydra/monitoring'", func() {
			So(res.String(), ShouldEqual, "https://smirnov.wiki/project/hydra/monitoring")
		})
	})
	Convey("Given 'sp/random'", t, func() {
		c, _ := parseDummyYaml()
		l := tokenize("sp/random")
		var res bytes.Buffer
		res.WriteString("https:/")

		expand(c, l.Front(), &res)

		Convey("result should equal 'https://smirnov.wiki/project/random'", func() {
			So(res.String(), ShouldEqual, "https://smirnov.wiki/project/random")
		})
	})
	Convey("Given 'sp/h/very/deep/path'", t, func() {
		c, _ := parseDummyYaml()
		l := tokenize("sp/h/very/deep/path")
		var res bytes.Buffer
		res.WriteString("https:/")

		expand(c, l.Front(), &res)

		Convey("result should equal 'https://smirnov.wiki/project/hydra/very/deep/path'", func() {
			So(res.String(), ShouldEqual, "https://smirnov.wiki/project/hydra/very/deep/path")
		})
	})
}
