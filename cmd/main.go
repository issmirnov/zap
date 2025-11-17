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
		log.Fatalf("Error parsing config file '%s'. Please fix syntax: %s\n", *configName, err)
	}

	// Perform extended validation of config.
	if *validate {
		if err := zap.ValidateConfig(c); err != nil {
			fmt.Printf("Configuration validation failed:\n%s\n", err.Error())
			os.Exit(1)
		}
		fmt.Println("Configuration validation successful - no errors detected.")
		os.Exit(0)
	}

	// Validate config before starting server
	if err := zap.ValidateConfig(c); err != nil {
		log.Fatalf("Configuration validation failed. Please fix errors before starting server:\n%s\n", err.Error())
	}

	context := &zap.Context{Config: c, Advertise: *advertise}

	// Try to update hosts file, but don't fail if we can't
	if err := zap.UpdateHosts(context); err != nil {
		log.Printf("Warning: Failed to update /etc/hosts file: %v", err)
		log.Println("Server will continue running, but DNS shortcuts may not work")
	}

	// Enable hot reload.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to create file watcher: %v", err)
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
		log.Fatalf("Failed to watch config directory: %v", err)
	}

	// Set up routes.
	router := SetupRouter(context)

	// Start the server
	serverAddr := fmt.Sprintf("%s:%d", *host, *port)
	fmt.Printf("Launching %s on %s\n", appName, serverAddr)
	fmt.Printf("Configuration file: %s\n", *configName)
	fmt.Printf("Health check: http://%s/healthz\n", serverAddr)
	fmt.Printf("Configuration view: http://%s/varz\n", serverAddr)

	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
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
