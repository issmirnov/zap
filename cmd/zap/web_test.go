package zap

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jeffail/gabs/v2"
	"github.com/ghodss/yaml"
	. "github.com/smartystreets/goconvey/convey"
)

// TODO: add tests that use erroneous Config.
// Will likely require injecting custom logger and intercepting error msgs.

// See https://elithrar.github.io/article/testing-http-handlers-go/ for comments.
func TestIndexHandler(t *testing.T) {
	Convey("Given app is set up with default Config", t, func() {
		c, err := loadTestYaml()
		So(err, ShouldBeNil)
		context := &Context{Config: c}
		appHandler := &CtxWrapper{context, IndexHandler}
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

			Convey("The result should be a 302 to https://github.com/issmirnov/zap/", func() {
				So(rr.Code, ShouldEqual, http.StatusFound)
				So(rr.Header().Get("Location"), ShouldEqual, "https://github.com/issmirnov/zap/")
			})
		})
		Convey("When we GET http://g/z/ with 'X-Forwarded-Host' set", func() {
			req, err := http.NewRequest("GET", "/z/", nil)
			So(err, ShouldBeNil)
			req.Header = map[string][]string{"X-Forwarded-Host": {"g"}}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			Convey("The result should be a 302 to https://github.com/issmirnov/zap/", func() {
				So(rr.Code, ShouldEqual, http.StatusFound)
				So(rr.Header().Get("Location"), ShouldEqual, "https://github.com/issmirnov/zap/")
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

			Convey("The result should be a 302 to https://github.com/issmirnov/zap/very/deep/path/", func() {
				So(rr.Code, ShouldEqual, http.StatusFound)
				So(rr.Header().Get("Location"), ShouldEqual, "https://github.com/issmirnov/zap/very/deep/path/")
			})
		})
		Convey("When we GET http://g/", func() {
			req, err := http.NewRequest("GET", "/", nil)
			So(err, ShouldBeNil)
			req.Host = "g"

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			Convey("The result should be a 302 to https://github.com/", func() {
				So(rr.Code, ShouldEqual, http.StatusFound)
				So(rr.Header().Get("Location"), ShouldEqual, "https://github.com/")
			})
		})
		Convey("When we GET http://fake/path", func() {
			req, err := http.NewRequest("GET", "/path", nil)
			So(err, ShouldBeNil)
			req.Host = "fake"

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			Convey("The result should be a 404", func() {
				So(rr.Code, ShouldEqual, http.StatusNotFound)
			})
		})
		Convey("When we GET http://g/s/", func() {
			req, err := http.NewRequest("GET", "/s/", nil)
			So(err, ShouldBeNil)
			req.Host = "g"

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			Convey("The result should be a 302 to https://github.com/search?q=", func() {
				So(rr.Code, ShouldEqual, http.StatusFound)
				So(rr.Header().Get("Location"), ShouldEqual, "https://github.com/search?q=")
			})
		})
		Convey("When we GET http://g/s/foo", func() {
			req, err := http.NewRequest("GET", "/s/foo", nil)
			So(err, ShouldBeNil)
			req.Host = "g"

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			Convey("The result should be a 302 to https://github.com/search?q=foo", func() {
				So(rr.Code, ShouldEqual, http.StatusFound)
				So(rr.Header().Get("Location"), ShouldEqual, "https://github.com/search?q=foo")
			})
		})
		Convey("When we GET http://g/s", func() {
			req, err := http.NewRequest("GET", "/s", nil)
			So(err, ShouldBeNil)
			req.Host = "g"

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			Convey("The result should be a 302 to https://github.com/search?q=", func() {
				So(rr.Code, ShouldEqual, http.StatusFound)
				So(rr.Header().Get("Location"), ShouldEqual, "https://github.com/search?q=")
			})
		})
		Convey("When we GET http://z/ with ssl_off", func() {
			req, err := http.NewRequest("GET", "/", nil)
			So(err, ShouldBeNil)
			req.Host = "z"

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			Convey("The result should be a 302 to http://zero.com/", func() {
				So(rr.Code, ShouldEqual, http.StatusFound)
				So(rr.Header().Get("Location"), ShouldEqual, "http://zero.com/")
			})
		})
		Convey("When we GET http://zz/ with ssl_off: no ", func() {
			req, err := http.NewRequest("GET", "/", nil)
			So(err, ShouldBeNil)
			req.Host = "zz"

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			Convey("The result should be a 302 to https://zero.ssl.on.com", func() {
				So(rr.Code, ShouldEqual, http.StatusFound)
				So(rr.Header().Get("Location"), ShouldEqual, "https://zero.ssl.on.com/")
			})
		})

		Convey("When we GET http://l/a with ssl_off ", func() {
			req, err := http.NewRequest("GET", "/a", nil)
			So(err, ShouldBeNil)
			req.Host = "l"

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			Convey("The result should be a 302 to http://localhost:8080", func() {
				So(rr.Code, ShouldEqual, http.StatusFound)
				So(rr.Header().Get("Location"), ShouldEqual, "http://localhost:8080")
			})
		})

		Convey("When we GET http://l/a/ with ssl_off ", func() {
			req, err := http.NewRequest("GET", "/a/", nil)
			So(err, ShouldBeNil)
			req.Host = "l"

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			Convey("The result should be a 302 to http://localhost:8080/", func() {
				So(rr.Code, ShouldEqual, http.StatusFound)
				So(rr.Header().Get("Location"), ShouldEqual, "http://localhost:8080/")
			})
		})
		Convey("When we GET http://l/a/s with ssl_off", func() {
			req, err := http.NewRequest("GET", "/a/s", nil)
			So(err, ShouldBeNil)
			req.Host = "l"

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			Convey("The result should be a 302 to http://localhost:8080/service", func() {
				So(rr.Code, ShouldEqual, http.StatusFound)
				So(rr.Header().Get("Location"), ShouldEqual, "http://localhost:8080/service")
			})
		})
		Convey("When we GET http://l/a/s/ with ssl_off", func() {
			req, err := http.NewRequest("GET", "/a/s/", nil)
			So(err, ShouldBeNil)
			req.Host = "l"

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			Convey("The result should be a 302 to http://localhost:8080/service/", func() {
				So(rr.Code, ShouldEqual, http.StatusFound)
				So(rr.Header().Get("Location"), ShouldEqual, "http://localhost:8080/service/")
			})
		})

		Convey("When we GET http://ch/ with schema set to 'chrome' ", func() {
			req, err := http.NewRequest("GET", "/", nil)
			So(err, ShouldBeNil)
			req.Host = "ch"

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			expected := "chrome://"
			Convey(fmt.Sprintf("The result should be a 302 to %s", expected), func() {
				So(rr.Code, ShouldEqual, http.StatusFound)
				So(rr.Header().Get("Location"), ShouldEqual, expected)
			})
		})

		Convey("When we GET http://ch/foobar with schema set to 'chrome' where 'foobar' isn't in the Config ", func() {
			req, err := http.NewRequest("GET", "/foobar", nil)
			So(err, ShouldBeNil)
			req.Host = "ch"

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			expected := "chrome://foobar"
			Convey(fmt.Sprintf("The result should be a 302 to %s", expected), func() {
				So(rr.Code, ShouldEqual, http.StatusFound)
				So(rr.Header().Get("Location"), ShouldEqual, expected)
			})
		})

		Convey("When we GET http://ch/v with schema set to 'chrome' ", func() {
			req, err := http.NewRequest("GET", "/v", nil)
			So(err, ShouldBeNil)
			req.Host = "ch"

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			expected := "chrome://version"
			Convey(fmt.Sprintf("The result should be a 302 to %s", expected), func() {
				So(rr.Code, ShouldEqual, http.StatusFound)
				So(rr.Header().Get("Location"), ShouldEqual, expected)
			})
		})

		Convey("When we GET http://ch/n/d with schema set to 'chrome' ", func() {
			req, err := http.NewRequest("GET", "/n/d", nil)
			So(err, ShouldBeNil)
			req.Host = "ch"

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			expected := "chrome://net-internals/#dns"
			Convey(fmt.Sprintf("The result should be a 302 to %s", expected), func() {
				So(rr.Code, ShouldEqual, http.StatusFound)
				So(rr.Header().Get("Location"), ShouldEqual, expected)
			})
		})

	})
}

func TestIndexHandlerErrorHandling(t *testing.T) {
	// Test case: nil configuration
	ctx := &Context{Config: nil}
	req := httptest.NewRequest("GET", "/", nil)
	req.Host = "test.com"
	w := httptest.NewRecorder()

	status, err := IndexHandler(ctx, w, req)
	if status != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, status)
	}
	if err == nil {
		t.Error("Expected error for nil configuration")
	}

	// Test case: host not found in configuration
	config := gabs.New()
	ctx = &Context{Config: config}
	req = httptest.NewRequest("GET", "/", nil)
	req.Host = "nonexistent.com"
	w = httptest.NewRecorder()

	status, err = IndexHandler(ctx, w, req)
	if status != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, status)
	}
	if err == nil {
		t.Error("Expected error for nonexistent host")
	}
}

// BenchmarkIndexHandler tests request processing speed when Context is preloaded.
// Run with go test -run=BenchmarkIndexHandler -bench=. // results: 500000x	2555 ns/op
func BenchmarkIndexHandler(b *testing.B) {
	c, _ := loadTestYaml()
	context := &Context{Config: c}
	appHandler := &CtxWrapper{context, IndexHandler}
	handler := http.Handler(appHandler)
	req, _ := http.NewRequest("GET", "/z", nil)
	req.Host = "g"
	rr := httptest.NewRecorder()
	for n := 0; n < b.N; n++ {
		handler.ServeHTTP(rr, req)
	}
}

func TestHealthCheckHandler(t *testing.T) {
	Convey("When we GET /healthz", t, func() {
		req, err := http.NewRequest("GET", "/healthz", nil)
		So(err, ShouldBeNil)
		req.Host = "sd"

		// We create a ResponseWriter (which satisfies http.ResponseWriter) to record the response.
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
	Convey("Given app is set up with default Config", t, func() {
		c, err := loadTestYaml()
		So(err, ShouldBeNil)
		context := &Context{Config: c}

		appHandler := &CtxWrapper{context, VarsHandler}
		handler := http.Handler(appHandler)
		Convey("When we GET /varz", func() {
			req, err := http.NewRequest("GET", "/varz", nil)
			So(err, ShouldBeNil)
			req.Host = "sd"

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			Convey("We should get a 200", func() {
				So(rr.Code, ShouldEqual, http.StatusOK)
			})
			Convey("It should be valid json", func() {
				_, err := yaml.YAMLToJSON(rr.Body.Bytes())
				So(err, ShouldBeNil)
			})
			Convey("It should equal the Config file", func() {
				conf, err := yaml.YAMLToJSON(c.Bytes())
				So(err, ShouldBeNil)

				resp, err := yaml.YAMLToJSON(rr.Body.Bytes())
				So(err, ShouldBeNil)

				// This does not work: "So(resp, ShouldEqual, []byte(jsonPrettyPrint(string(conf))))"
				// We get a nicely formatted response, but when we feed it into YAMLToJSON it collapses our nice
				// newlines. As a result, directly comparing the byte arrays here is a nogo. Therefore, we cheat
				// and utilize the separately tested jsonPrettyPrint to idempotently indent the JSON and compare that.
				So(jsonPrettyPrint(string(resp)), ShouldEqual, jsonPrettyPrint(string(conf)))
			})
		})
	})
}
