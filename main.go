package main

import (
    "flag"
    "fmt"
    "net/http"
    "os"
)

const appName = "zap"

var version = "develop"

func main() {

    var (
        configName = flag.String("config", "c.yml", "config file")
        port       = flag.Int("port", 8927, "port to bind to")
        host       = flag.String("host", "127.0.0.1", "host interface")
        i          = flag.String("index", "start", "string to append if path has trailing slash")
        v          = flag.Bool("v", false, "print version info")
    )
    flag.Parse()

    if *v {
        fmt.Println(version)
        os.Exit(0)
    }

    c, err := parseYaml(*configName)
    if err != nil {
        fmt.Printf("error: %s\n", err)
        return
    }
    context := &context{config: c, index: *i}

    http.Handle("/", ctxWrapper{context, IndexHandler})
    http.Handle("/varz", ctxWrapper{context, VarsHandler})
    http.HandleFunc("/healthz", HealthHandler)

    // TODO check for errors - addr in use, sudo issues, etc.
    fmt.Printf("Launching %s on %s:%d\n", appName, *host, *port)
    http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil)
}
