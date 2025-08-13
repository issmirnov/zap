package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/issmirnov/zap/cmd/zap"

	"github.com/fsnotify/fsnotify"

	"github.com/julienschmidt/httprouter"
)

const appName = "zap"

// Used in version printer, set by GoReleaser.
var version = "develop"

func main() {
	var (
		configName = flag.String("config", "c.yml", "config file")
		port       = flag.Int("port", 8927, "port to bind to")
		host       = flag.String("host", "127.0.0.1", "host address to bind to")
		advertise  = flag.String("advertise", "127.0.0.1", "IP to advertise, used in /etc/hosts")
		v          = flag.Bool("v", false, "print version info")
		validate   = flag.Bool("validate", false, "load config file and check for errors")
	)
	flag.Parse()

	if *v {
		fmt.Println(version)
		os.Exit(0)
	}

	// load config for first time.
	c, err := zap.ParseYaml(*configName)
	if err != nil {
		log.Printf("Error parsing config file. Please fix syntax: %s\n", err)
		return
	}

	// Perform extended validation of config.
	if *validate {
		if err := zap.ValidateConfig(c); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println("No errors detected.")
		os.Exit(0)
	}

	context := &zap.Context{Config: c, Advertise: *advertise}
	zap.UpdateHosts(context) // sync changes since last run.

	// Enable hot reload.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := watcher.Close(); err != nil {
			log.Printf("Error closing watcher: %v", err)
		}
	}()

	cb := zap.MakeReloadCallback(context, *configName)
	go zap.WatchConfigFileChanges(watcher, *configName, cb)
	err = watcher.Add(path.Dir(*configName))
	if err != nil {
		log.Fatal(err)
	}

	// Set up routes.
	router := SetupRouter(context)

	// TODO check for errors - addr in use, sudo issues, etc.
	fmt.Printf("Launching %s on %s:%d\n", appName, *host, *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), router))
}

func SetupRouter(context *zap.Context) *httprouter.Router {
	router := httprouter.New()
	router.Handler("GET", "/", zap.CtxWrapper{Context: context, H: zap.IndexHandler})
	router.Handler("GET", "/varz", zap.CtxWrapper{Context: context, H: zap.VarsHandler})
	router.HandlerFunc("GET", "/healthz", zap.HealthHandler)

	// https://github.com/julienschmidt/httprouter is having issues with
	// wildcard handling. As a result, we have to register index handler
	// as the fallback. Fix incoming.
	router.NotFound = zap.CtxWrapper{Context: context, H: zap.IndexHandler}
	return router
}
