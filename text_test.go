package main

import (
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTokenizer(t *testing.T) {
	Convey("Given a string 'g/z'", t, func() {
		l := tokenize("g/z")
		Convey("The resulting list should", func() {

			Convey("Have length 2", func() {
				So(l.Len(), ShouldEqual, 2)
			})
			Convey("The first element should be equal to 'g'", func() {
				So(l.Front().Value, ShouldEqual, "g")
			})
			Convey("The last element should be equal to 'z'", func() {
				So(l.Back().Value, ShouldEqual, "z")
			})
		})
	})

	Convey("Given a string 'e/a/extratext'", t, func() {
		l := tokenize("e/a/extratext")
		Convey("The resulting list should have length 3", func() {
			So(l.Len(), ShouldEqual, 3)
		})
	})

	Convey("Given a string 'e/a/extratext/'", t, func() {
		l := tokenize("e/a/extratext/")
		Convey("The resulting list should have length 4", func() {
			So(l.Len(), ShouldEqual, 4) // Since we have nil terminator.
		})
	})
}

func TestExpander(t *testing.T) {

	Convey("Given 'g/z'", t, func() {
		c, _ := loadTestYaml()
		l := tokenize("g/z")
		var res bytes.Buffer
		res.WriteString(httpsPrefix)

		expandPath(c, l.Front(), &res)

		Convey("result should equal 'https://github.com/issmirnov/zap'", func() {
			So(res.String(), ShouldEqual, "https://github.com/issmirnov/zap")
		})
	})
	// Convey("Given 'e/n'", t, func() {
	//     c, _ := loadTestYaml()
	//     l := tokenize("e/n ")
	//     var res bytes.Buffer
	//     res.WriteString(httpsPrefix)
	//
	//     expandPath(c, l.Front(), &res)
	//
	//     Convey("result should equal 'https://example.com/999'", func() {
	//         So(res.String(), ShouldEqual, "https://example.com/999")
	//     })
	// })
	Convey("Given 'g/z/extratext'", t, func() {
		c, _ := loadTestYaml()
		l := tokenize("g/z/extratext")
		var res bytes.Buffer
		res.WriteString(httpsPrefix)

		expandPath(c, l.Front(), &res)

		Convey("result should equal 'https://github.com/issmirnov/zap/extratext'", func() {
			So(res.String(), ShouldEqual, "https://github.com/issmirnov/zap/extratext")
		})
	})
	Convey("Given 'g/'", t, func() {
		c, _ := loadTestYaml()
		l := tokenize("g/")
		var res bytes.Buffer
		res.WriteString(httpsPrefix)

		expandPath(c, l.Front(), &res)

		Convey("result should equal 'https://github.com/'", func() {
			So(res.String(), ShouldEqual, "https://github.com/")
		})
	})
	Convey("Given 'g/z/very/deep/path'", t, func() {
		c, _ := loadTestYaml()
		l := tokenize("g/z/very/deep/path")
		var res bytes.Buffer
		res.WriteString(httpsPrefix)

		expandPath(c, l.Front(), &res)

		Convey("result should equal 'https://github.com/issmirnov/zap/very/deep/path'", func() {
			So(res.String(), ShouldEqual, "https://github.com/issmirnov/zap/very/deep/path")
		})
	})
	Convey("Given 'g/s/foobar'", t, func() {
		c, _ := loadTestYaml()
		l := tokenize("g/s/foobar")
		var res bytes.Buffer
		res.WriteString(httpsPrefix)

		expandPath(c, l.Front(), &res)

		Convey("result should equal 'https://github.com/search?q=foobar'", func() {
			So(res.String(), ShouldEqual, "https://github.com/search?q=foobar")
		})
	})
	Convey("Given 'g/s/foo/bar/baz'", t, func() {
		c, _ := loadTestYaml()
		l := tokenize("g/s/foo/bar/baz")
		var res bytes.Buffer
		res.WriteString(httpsPrefix)

		expandPath(c, l.Front(), &res)

		Convey("result should equal 'https://github.com/search?q=foo/bar/baz'", func() {
			So(res.String(), ShouldEqual, "https://github.com/search?q=foo/bar/baz")
		})
	})
	Convey("Given 'g/s/foo/bar/baz/'", t, func() {
		c, _ := loadTestYaml()
		l := tokenize("g/s/foo/bar/baz/")
		var res bytes.Buffer
		res.WriteString(httpsPrefix)

		expandPath(c, l.Front(), &res)

		Convey("result should equal 'https://github.com/search?q=foo/bar/baz/'", func() {
			So(res.String(), ShouldEqual, "https://github.com/search?q=foo/bar/baz/")
		})
	})
}
