package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/render"
	"gopkg.in/yaml.v2"

	"github.com/nate-martinez/faux-twitter/server/log"
	"github.com/nate-martinez/faux-twitter/server/util"
)

var (
	mut sync.RWMutex

	global = &Config{}
)

type Config struct {
	Source string `yaml:"source" json:"source"`

	LogLevel string `yaml:"logLevel" json:"log_level"`
}

func LoadFile(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("error reading config from '%s': %s", path, err)
		log.Log.Error(err)
		return err
	}
	conf := &Config{}
	if err = yaml.Unmarshal(raw, conf); err != nil {
		err = fmt.Errorf("error unmarshalling config at '%s': %s", path, err)
		log.Log.Error(err)
		return err
	}
	conf.SetDefaults(path, true)
	conf.load()
	global.Print()
	return nil
}

func GetConfig(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, global)
}

func PostConfig(w http.ResponseWriter, r *http.Request) {
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		render.JSON(w, r, util.StatusResponse{Status: "error", Message: fmt.Sprintf("error reading request body: %s", err)})
		return
	}
	conf := &Config{}
	if err = json.Unmarshal(raw, conf); err != nil {
		w.WriteHeader(400)
		render.JSON(w, r, util.StatusResponse{Status: "error", Message: fmt.Sprintf("error unmarshaling config: %s", err)})
		return
	}
	conf.SetDefaults(r.RemoteAddr, false)
	conf.load()
	render.JSON(w, r, util.StatusResponse{Status: "success"})
}

func (c *Config) Print() {
	log.Log.Infof("LogLevel: %s", c.LogLevel)
}

const configStr = "{Source: %s, LogLevel: %s}"

func (c *Config) String() string {
	return fmt.Sprintf(configStr,
		c.Source,
		c.LogLevel,
	)
}

const (
	LogLevelDefault = "Info"
)

func (c *Config) SetDefaults(source string, shouldLog bool) {
	c.Source = source
	switch c.LogLevel {
	case "Error", "Warn", "Info", "Debug":
	default:
		if shouldLog {
			log.Log.Infof("setting LogLevel to default %s", LogLevelDefault)
			c.LogLevel = LogLevelDefault
		}
	}
}

func (c *Config) load() {
	mut.Lock()
	defer mut.Unlock()
	log.Log.Infof("loading config %s from %s...", c, c.Source)

	// manual hooks
	if c.LogLevel != global.LogLevel {
		log.SetLogLevel(c.LogLevel)
	}

	global = c
	log.Log.Infof("finished loading config from %s", c.Source)
}

// getter methods

func LogLevel() string {
	mut.RLock()
	defer mut.RUnlock()
	return global.LogLevel
}
