package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"

	"github.com/Jeffail/gabs"
	"github.com/fsnotify/fsnotify"
	"github.com/ghodss/yaml"
)

// expandKey is used to get expand mapping
const expandKey string = "expand"

// queryKey is used to get search prefix
const queryKey string = "query"

// parseYaml takes a file name and returns a gabs config object.
func parseYaml(fname string) (*gabs.Container, error) {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		fmt.Printf("Unable to read file: %s\n", err.Error())
		return nil, err
	}
	d, jsonErr := yaml.YAMLToJSON([]byte(data))
	if jsonErr != nil {
		fmt.Printf("Error encoding input to JSON.\n%s\n", jsonErr.Error())
		return nil, jsonErr
	}
	j, _ := gabs.ParseJSON(d)
	return j, nil
}

func watchChanges(watcher *fsnotify.Watcher, fname string, cb func()) {
	for {
		select {
		case event := <-watcher.Events:
			// You may wonder why we can't just listen for "Write" events. The reason is that vim (and other editors)
			// will create swap files, and when you write they delete the original and rename the swap file. This is great
			// for resolving system crashes, but also completely incompatible with inotify and other fswatch implementations.
			// Thus, we check that the file of interest might be created as well.
			if (event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write) && event.Name == fname {
				cb()
			}
		case e := <-watcher.Errors:
			log.Println("error:", e)
		}
	}
}

// TODO: add tests. simulate touching a file.
// updateHosts will attempt to write the zap list of shortcuts
// to /etc/hosts. It will gracefully fail if there are not enough
// permissions to do so.
func updateHosts(c *context) {
	hostPath := "/etc/hosts"

	// 1. read file, prep buffer.
	data, err := ioutil.ReadFile(hostPath)
	if err != nil {
		log.Println("open config: ", err)
	}
	var replacement bytes.Buffer

	// 2. generate payload.
	replacement.WriteString("### Zap Shortcuts :start ##\n")
	children, _ := c.config.ChildrenMap()
	for k := range children {
		replacement.WriteString(fmt.Sprintf("127.0.0.1 %s\n", k))
	}
	replacement.WriteString("### Zap Shortcuts :end ##")

	// 3. run regexp, prepare file content for write.
	zapBlock := regexp.MustCompile("(###(.*)##)\n(.|\n)*(###(.*)##)")
	updatedFile := zapBlock.ReplaceAllString(string(data), replacement.String())

	// If no changes, we have no presence yet. Append our data.
	if updatedFile == string(data) {
		updatedFile += replacement.String()
	}

	// 4. Attempt write to file.
	err = ioutil.WriteFile(hostPath, []byte(updatedFile), 0644)
	if err != nil {
		log.Printf("Error writing to '%s': %s\n", hostPath, err.Error())
	}

}

// makeCallback returns a func that that updates global state.
func makeCallback(c *context, configName string) func() {
	return func() {
		data, err := parseYaml(configName)
		if err != nil {
			log.Printf("Error in new config: %s. Fallback to old config.", err)
			return
		}

		// Update config atomically
		c.configMtx.Lock()
		c.config = data
		c.configMtx.Unlock()

		// Sync DNS entries.
		updateHosts(c)
		return
	}
}
