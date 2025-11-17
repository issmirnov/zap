package zap

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/fsnotify/fsnotify"
	"github.com/ghodss/yaml"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/afero"
)

const (
	delimStart  = "### Zap Shortcuts :start ##\n"
	delimEnd    = "### Zap Shortcuts :end ##\n"
	expandKey   = "expand"
	queryKey    = "query"
	portKey     = "port"
	passKey     = "*"
	sslKey      = "ssl_off"
	schemaKey   = "schema"
	httpsPrefix = "https:/" // second slash appended in expandPath() call
	httpPrefix  = "http:/"  // second slash appended in expandPath() call
)

// Sentinel value used to indicate set membership.
var exists = struct{}{}

// Afero is a filesystem wrapper providing util methods
// and easy test mocks.
var Afero = &afero.Afero{Fs: afero.NewOsFs()}

// parseYamlString takes a raw string and attempts to load it.
func parseYamlString(config string) (*gabs.Container, error) {
	if config == "" {
		return nil, fmt.Errorf("empty configuration string provided")
	}

	d, jsonErr := yaml.YAMLToJSON([]byte(config))
	if jsonErr != nil {
		return nil, fmt.Errorf("failed to parse YAML configuration: %w", jsonErr)
	}

	j, err := gabs.ParseJSON(d)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON configuration: %w", err)
	}
	return j, nil
}

// ParseYaml takes a file name and returns a gabs Config object.
func ParseYaml(fname string) (*gabs.Container, error) {
	if fname == "" {
		return nil, fmt.Errorf("no configuration file specified")
	}

	data, err := Afero.ReadFile(fname)
	if err != nil {
		return nil, fmt.Errorf("unable to read configuration file '%s': %w", fname, err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("configuration file '%s' is empty", fname)
	}

	d, jsonErr := yaml.YAMLToJSON(data)
	if jsonErr != nil {
		return nil, fmt.Errorf("failed to parse YAML configuration in file '%s': %w", fname, jsonErr)
	}

	j, err := gabs.ParseJSON(d)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON configuration in file '%s': %w", fname, err)
	}
	return j, nil
}

// ValidateConfig verifies that there are no unexpected values in the Config file.
// at each level of the Config, we should either have a KV for expansions, or a leaf node
// with the values oneof "expand", "query", "ssl_off" that map to a string.
func ValidateConfig(c *gabs.Container) error {
	if c == nil {
		return fmt.Errorf("configuration is nil")
	}

	var errors *multierror.Error
	children := c.ChildrenMap()
	seenKeys := make(map[string]struct{})

	for k, v := range children {
		// Check if key already seen
		if _, ok := seenKeys[k]; ok {
			errors = multierror.Append(errors, fmt.Errorf("duplicate key '%s' detected in configuration", k))
		} else {
			seenKeys[k] = exists // mark key as seen
		}

		// Validate all children
		switch k {
		case expandKey:
			// check that v is a string or number, else return error.
			if _, ok := v.Data().(string); !ok {
				if _, ok := v.Data().(float64); !ok {
					errors = multierror.Append(errors, fmt.Errorf("expected string or number value for 'expand' key, got: %T (%v)", v.Data(), v.Data()))
				}
			}
		case queryKey:
			// check that v is a string, else return error.
			if _, ok := v.Data().(string); !ok {
				errors = multierror.Append(errors, fmt.Errorf("expected string value for 'query' key, got: %T (%v)", v.Data(), v.Data()))
			}
		case schemaKey:
			// check that v is a string, else return error.
			if _, ok := v.Data().(string); !ok {
				errors = multierror.Append(errors, fmt.Errorf("expected string value for 'schema' key, got: %T (%v)", v.Data(), v.Data()))
			}
		case portKey:
			// check that v is a float64, else return error.
			if _, ok := v.Data().(float64); !ok {
				errors = multierror.Append(errors, fmt.Errorf("expected number value for 'port' key, got: %T (%v)", v.Data(), v.Data()))
			}
		case sslKey:
			// check that v is a boolean, else return error.
			if _, ok := v.Data().(bool); !ok {
				errors = multierror.Append(errors, fmt.Errorf("expected boolean value for 'ssl_off' key, got: %T (%v)", v.Data(), v.Data()))
			}
		default:
			// Check if we have an unknown string here.
			if _, ok := v.Data().(string); ok {
				errors = multierror.Append(errors, fmt.Errorf("unexpected string value '%v' under key '%s' - this key is not recognized", v.Data(), k))
			}
			// recurse, collect any errors.
			if err := ValidateConfig(v); err != nil {
				errors = multierror.Append(errors, fmt.Errorf("validation error in section '%s': %w", k, err))
			}
		}
	}
	return errors.ErrorOrNil()
}

// WatchConfigFileChanges will attach an fsnotify watcher to the config file, and trigger
// the cb function when the file is updated.
func WatchConfigFileChanges(watcher *fsnotify.Watcher, fname string, cb func()) {
	for {
		select {
		case event := <-watcher.Events:
			// You may wonder why we can't just listen for "Write" events. The reason is that vim (and other editors)
			// will create swap files, and when you write they delete the original and rename the swap file. This is great
			// for resolving system crashes, but also completely incompatible with inotify and other fswatch implementations.
			// Thus, we check that the file of interest might be created as well.
			updated := event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write
			zapconf := filepath.Clean(event.Name) == fname
			if updated && zapconf {
				log.Printf("Configuration file '%s' changed, reloading...", fname)
				cb()
				log.Printf("Configuration reloaded successfully")
			}
		case e := <-watcher.Errors:
			log.Printf("File watcher error: %v", e)
		}
	}
}

// TODO: add tests. simulate touching a file.
// UpdateHosts will attempt to write the zap list of shortcuts
// to /etc/hosts. It will gracefully fail if there are not enough
// permissions to do so.
func UpdateHosts(c *Context) error {
	hostPath := "/etc/hosts"

	// 1. read file, prep buffer.
	data, err := os.ReadFile(hostPath)
	if err != nil {
		return fmt.Errorf("failed to read hosts file '%s': %w", hostPath, err)
	}
	var replacement bytes.Buffer

	// 2. generate payload.
	replacement.WriteString(delimStart)
	children := c.Config.ChildrenMap()
	for k := range children {
		replacement.WriteString(fmt.Sprintf("%s %s\n", c.Advertise, k))
	}
	replacement.WriteString(delimEnd)

	// 3. Generate new file content
	var updatedFile string
	if !strings.Contains(string(data), delimStart) {
		updatedFile = string(data) + replacement.String()
	} else {
		zapBlock := regexp.MustCompile("(###(.*)##)\n(.|\n)*(###(.*)##\n)")
		updatedFile = zapBlock.ReplaceAllString(string(data), replacement.String())
	}

	// 4. Attempt write to file.
	err = os.WriteFile(hostPath, []byte(updatedFile), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to hosts file '%s': %w", hostPath, err)
	}

	return nil
}

// MakeReloadCallback returns a func that that reads the config file and updates global state.
func MakeReloadCallback(c *Context, configName string) func() {
	return func() {
		data, err := ParseYaml(configName)
		if err != nil {
			log.Printf("Error loading new Config: %s. Fallback to old Config.", err)
			return
		}
		err = ValidateConfig(data)
		if err != nil {
			log.Printf("Error validating new Config: %s. Fallback to old Config.", err)
			return
		}

		// Update Config atomically
		c.ConfigMtx.Lock()
		c.Config = data
		c.ConfigMtx.Unlock()

		// Sync DNS entries.
		if err := UpdateHosts(c); err != nil {
			log.Printf("Warning: Failed to update hosts file during reload: %v", err)
		}
	}
}
