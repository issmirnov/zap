package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ghodss/yaml"
	. "github.com/smartystreets/goconvey/convey"
)

// See https://elithrar.github.io/article/testing-http-handlers-go/ for comments.
func TestIndexHandler(t *testing.T) {
	Convey("Given app is set up with default config", t, func() {
		c, err := parseDummyYaml()
		So(err, ShouldBeNil)
		context := &context{config: c, index: "start"}
		appHandler := &ctxWrapper{context, IndexHandler}
		handler := http.Handler(appHandler)
		Convey("When we GET http://g/z", func() {
			req, err := http.NewRequest("GET", "/z", nil)
			So(err, ShouldBeNil)
			req.Host = "g"

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			Convey("The result should be a 302 to https://github.com/issmirnov/zap", func() {
				So(rr.Code, ShouldEqual, http.StatusFound)
				So(rr.Header().Get("Location"), ShouldEqual, "https://github.com/issmirnov/zap")
			})
		})
		Convey("When we GET http://g/z/", func() {
			req, err := http.NewRequest("GET", "/z/", nil)
			So(err, ShouldBeNil)
			req.Host = "g"

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			Convey("The result should be a 302 to https://github.com/issmirnov/zap/start", func() {
				So(rr.Code, ShouldEqual, http.StatusFound)
				So(rr.Header().Get("Location"), ShouldEqual, "https://github.com/issmirnov/zap/start")
			})
		})
		Convey("When we GET http://g/z/very/deep/path", func() {
			req, err := http.NewRequest("GET", "/z/very/deep/path", nil)
			So(err, ShouldBeNil)
			req.Host = "g"

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			Convey("The result should be a 302 to https://github.com/issmirnov/zap/very/deep/path", func() {
				So(rr.Code, ShouldEqual, http.StatusFound)
				So(rr.Header().Get("Location"), ShouldEqual, "https://github.com/issmirnov/zap/very/deep/path")
			})
		})
		Convey("When we GET http://g/z/very/deep/path/", func() {
			req, err := http.NewRequest("GET", "/z/very/deep/path/", nil)
			So(err, ShouldBeNil)
			req.Host = "g"

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			Convey("The result should be a 302 to https://github.com/issmirnov/zap/very/deep/path/start", func() {
				So(rr.Code, ShouldEqual, http.StatusFound)
				So(rr.Header().Get("Location"), ShouldEqual, "https://github.com/issmirnov/zap/very/deep/path/start")
			})
		})
		Convey("When we GET http://g/", func() {
			req, err := http.NewRequest("GET", "/", nil)
			So(err, ShouldBeNil)
			req.Host = "g"

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			Convey("The result should be a 302 to https://github.com/start", func() {
				So(rr.Code, ShouldEqual, http.StatusFound)
				So(rr.Header().Get("Location"), ShouldEqual, "https://github.com/start")
			})
		})
	})
	Convey("Given app is set up with non default config", t, func() {
		c, err := parseDummyYaml()
		So(err, ShouldBeNil)
		context := &context{config: c, index: "otherIndex"}
		appHandler := &ctxWrapper{context, IndexHandler}
		handler := http.Handler(appHandler)
		Convey("When we GET http://g/", func() {
			req, err := http.NewRequest("GET", "/", nil)
			So(err, ShouldBeNil)
			req.Host = "g"

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			Convey("The result should be a 302 to https://github.com/otherIndex", func() {
				So(rr.Code, ShouldEqual, http.StatusFound)
				So(rr.Header().Get("Location"), ShouldEqual, "https://github.com/otherIndex")
			})
		})
		Convey("When we GET http://g/z/very/deep/path", func() {
			req, err := http.NewRequest("GET", "/z/very/deep/path", nil)
			So(err, ShouldBeNil)
			req.Host = "g"

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			Convey("The result should be a 302 to https://github.com/issmirnov/zap/very/deep/path", func() {
				So(rr.Code, ShouldEqual, http.StatusFound)
				So(rr.Header().Get("Location"), ShouldEqual, "https://github.com/issmirnov/zap/very/deep/path")
			})
		})
		Convey("When we GET http://g/z/very/deep/path/", func() {
			req, err := http.NewRequest("GET", "/z/very/deep/path/", nil)
			So(err, ShouldBeNil)
			req.Host = "g"

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			Convey("The result should be a 302 to https://github.com/issmirnov/zap/very/deep/path/otherIndex", func() {
				So(rr.Code, ShouldEqual, http.StatusFound)
				So(rr.Header().Get("Location"), ShouldEqual, "https://github.com/issmirnov/zap/very/deep/path/otherIndex")
			})
		})
	})
}

// BenchmarkIndexHandler tests request processing geed when context is preloaded.
// Run with go test -run=BenchmarkIndexHandler -bench=. // results: 500000x	2555 ns/op
func BenchmarkIndexHandler(b *testing.B) {
	c, _ := parseDummyYaml()
	context := &context{config: c}
	appHandler := &ctxWrapper{context, IndexHandler}
	handler := http.Handler(appHandler)
	req, _ := http.NewRequest("GET", "/p", nil)
	req.Host = "sd"
	rr := httptest.NewRecorder()
	for n := 0; n < b.N; n++ {
		handler.ServeHTTP(rr, req)
	}
}

func TestHealthCheckHandler(t *testing.T) {
	Convey("When we GET /zealthz", t, func() {
		req, err := http.NewRequest("GET", "/zealthz", nil)
		So(err, ShouldBeNil)
		req.Host = "sd"

		// We create a RegonseRecorder (which satisfies http.RegonseWriter) to record the regonse.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(HealthHandler)
		handler.ServeHTTP(rr, req)

		Convey("We should get a 200", func() {
			So(rr.Code, ShouldEqual, http.StatusOK)
			So(rr.Body.String(), ShouldEqual, "OK")
		})
	})
}

func TestVarzHandler(t *testing.T) {
	Convey("Given app is set up with default config", t, func() {
		c, err := parseDummyYaml()
		So(err, ShouldBeNil)
		context := &context{config: c}

		appHandler := &ctxWrapper{context, VarsHandler}
		handler := http.Handler(appHandler)
		Convey("When we GET /varz", func() {
			req, err := http.NewRequest("GET", "/varz", nil)
			So(err, ShouldBeNil)
			req.Host = "sd"

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			Convey("We should get a 200", func() {
				So(rr.Code, ShouldEqual, http.StatusOK)
				conf, err := yaml.YAMLToJSON(c.Bytes())
				So(err, ShouldBeNil)
				d, err := yaml.YAMLToJSON(rr.Body.Bytes())
				So(err, ShouldBeNil)
				So(string(d), ShouldContainSubstring, string(conf))
			})
		})
	})
}
