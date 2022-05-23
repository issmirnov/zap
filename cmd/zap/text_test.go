package zap

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

		ExpandPath(c, l.Front(), &res)

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
	//     ExpandPath(c, l.Front(), &res)
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

		ExpandPath(c, l.Front(), &res)

		Convey("result should equal 'https://github.com/issmirnov/zap/extratext'", func() {
			So(res.String(), ShouldEqual, "https://github.com/issmirnov/zap/extratext")
		})
	})
	Convey("Given 'g/'", t, func() {
		c, _ := loadTestYaml()
		l := tokenize("g/")
		var res bytes.Buffer
		res.WriteString(httpsPrefix)

		ExpandPath(c, l.Front(), &res)

		Convey("result should equal 'https://github.com/'", func() {
			So(res.String(), ShouldEqual, "https://github.com/")
		})
	})
	Convey("Given 'g/z/very/deep/path'", t, func() {
		c, _ := loadTestYaml()
		l := tokenize("g/z/very/deep/path")
		var res bytes.Buffer
		res.WriteString(httpsPrefix)

		ExpandPath(c, l.Front(), &res)

		Convey("result should equal 'https://github.com/issmirnov/zap/very/deep/path'", func() {
			So(res.String(), ShouldEqual, "https://github.com/issmirnov/zap/very/deep/path")
		})
	})
	Convey("Given 'g/s/foobar'", t, func() {
		c, _ := loadTestYaml()
		l := tokenize("g/s/foobar")
		var res bytes.Buffer
		res.WriteString(httpsPrefix)

		ExpandPath(c, l.Front(), &res)

		Convey("result should equal 'https://github.com/search?q=foobar'", func() {
			So(res.String(), ShouldEqual, "https://github.com/search?q=foobar")
		})
	})
	Convey("Given 'g/s/foo/bar/baz'", t, func() {
		c, _ := loadTestYaml()
		l := tokenize("g/s/foo/bar/baz")
		var res bytes.Buffer
		res.WriteString(httpsPrefix)

		ExpandPath(c, l.Front(), &res)

		Convey("result should equal 'https://github.com/search?q=foo/bar/baz'", func() {
			So(res.String(), ShouldEqual, "https://github.com/search?q=foo/bar/baz")
		})
	})
	Convey("Given 'g/s/foo/bar/baz/'", t, func() {
		c, _ := loadTestYaml()
		l := tokenize("g/s/foo/bar/baz/")
		var res bytes.Buffer
		res.WriteString(httpsPrefix)

		ExpandPath(c, l.Front(), &res)

		Convey("result should equal 'https://github.com/search?q=foo/bar/baz/'", func() {
			So(res.String(), ShouldEqual, "https://github.com/search?q=foo/bar/baz/")
		})
	})
	Convey("Given 'g/query/homebrew'", t, func() {
		c, _ := loadTestYaml()
		l := tokenize("g/query/homebrew")
		var res bytes.Buffer
		res.WriteString(httpsPrefix)

		ExpandPath(c, l.Front(), &res)

		Convey("result should equal 'https://github.com/query/homebrew'", func() {
			So(res.String(), ShouldEqual, "https://github.com/query/homebrew")
		})
	})
	Convey("Given 'wc/1/*/3/four'", t, func() {
		c, _ := loadTestYaml()
		l := tokenize("wc/1/*/3/four")
		var res bytes.Buffer
		res.WriteString(httpsPrefix)

		ExpandPath(c, l.Front(), &res)

		Convey("result should equal 'https://wildcard.com/1/*/3/4'", func() {
			So(res.String(), ShouldEqual, "https://wildcard.com/1/*/3/4")
		})
	})
	Convey("Given 'wc/1/2/3/four'", t, func() {
		c, _ := loadTestYaml()
		l := tokenize("wc/1/2/3/four")
		var res bytes.Buffer
		res.WriteString(httpsPrefix)

		ExpandPath(c, l.Front(), &res)

		Convey("result should equal 'https://wildcard.com/1/2/3/4'", func() {
			So(res.String(), ShouldEqual, "https://wildcard.com/1/2/3/4")
		})
	})
	Convey("Given 'ak/hi'", t, func() {
		c, _ := loadTestYaml()
		l := tokenize("ak/hi")
		var res bytes.Buffer
		res.WriteString(httpsPrefix)

		ExpandPath(c, l.Front(), &res)

		Convey("result should equal 'https://kafka.apache.org/contact", func() {
			So(res.String(), ShouldEqual, "https://kafka.apache.org/contact")
		})
	})
	Convey("Given 'ak/23'", t, func() {
		c, _ := loadTestYaml()
		l := tokenize("ak/23")
		var res bytes.Buffer
		res.WriteString(httpsPrefix)

		ExpandPath(c, l.Front(), &res)

		Convey("result should equal 'https://kafka.apache.org/23", func() {
			So(res.String(), ShouldEqual, "https://kafka.apache.org/23")
		})
	})
	Convey("Given 'ak/23/j", t, func() {
		c, _ := loadTestYaml()
		l := tokenize("ak/23/j")
		var res bytes.Buffer
		res.WriteString(httpsPrefix)

		ExpandPath(c, l.Front(), &res)

		Convey("result should equal 'https://kafka.apache.org/23/javadoc/index.html?overview-summary.html", func() {
			So(res.String(), ShouldEqual, "https://kafka.apache.org/23/javadoc/index.html?overview-summary.html")
		})
	})
	Convey("Given 'ak/expand/j'", t, func() {
		c, _ := loadTestYaml()
		l := tokenize("ak/expand/j")
		var res bytes.Buffer
		res.WriteString(httpsPrefix)

		ExpandPath(c, l.Front(), &res)

		Convey("result should equal 'https://kafka.apache.org/expand/javadoc/index.html?overview-summary.html", func() {
			So(res.String(), ShouldEqual, "https://kafka.apache.org/expand/javadoc/index.html?overview-summary.html")
		})
	})
	Convey("Given 'g/s/me'", t, func() {
		c, _ := loadTestYaml()
		l := tokenize("g/s/me")
		var res bytes.Buffer
		res.WriteString(httpsPrefix)

		ExpandPath(c, l.Front(), &res)

		Convey("result should equal 'https://github.com/search?q=issmirnov'", func() {
			So(res.String(), ShouldEqual, "https://github.com/search?q=issmirnov")
		})
	})
	Convey("Given 'g/s/me/z'", t, func() {
		c, _ := loadTestYaml()
		l := tokenize("g/s/me/z")
		var res bytes.Buffer
		res.WriteString(httpsPrefix)

		ExpandPath(c, l.Front(), &res)

		Convey("result should equal 'https://github.com/search?q=issmirnov/zap'", func() {
			So(res.String(), ShouldEqual, "https://github.com/search?q=issmirnov/zap")
		})
	})
	Convey("Given 'g/s/ak'", t, func() {
		c, _ := loadTestYaml()
		l := tokenize("g/s/ak")
		var res bytes.Buffer
		res.WriteString(httpsPrefix)

		ExpandPath(c, l.Front(), &res)

		Convey("result should equal 'https://github.com/search?q=apache/kafka'", func() {
			So(res.String(), ShouldEqual, "https://github.com/search?q=apache/kafka")
		})
	})
	Convey("Given 'g/s/ak/c'", t, func() {
		c, _ := loadTestYaml()
		l := tokenize("g/s/ak/c")
		var res bytes.Buffer
		res.WriteString(httpsPrefix)

		ExpandPath(c, l.Front(), &res)

		Convey("result should equal 'https://github.com/search?q=apache/kafka+connect'", func() {
			So(res.String(), ShouldEqual, "https://github.com/search?q=apache/kafka+connect")
		})
	})
}
