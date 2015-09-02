package config

import (
	"github.com/spf13/cast"
	"gopkg.in/yaml.v2"
	"log"
	"reflect"
	"strings"
	"time"
)

var keyDelim string = "."

// RawConfig is essentially repository for configurations
type RawConfig map[string]interface{}

func (c RawConfig) searchMap(source map[string]interface{}, path []string) interface{} {

	if len(path) == 0 {
		return source
	}

	if next, ok := source[path[0]]; ok {
		switch next.(type) {
		case map[interface{}]interface{}:
			return c.searchMap(cast.ToStringMap(next), path[1:])
		case map[string]interface{}:
			// Type assertion is safe here since it is only reached
			// if the type of `next` is the same as the type being asserted
			return c.searchMap(next.(map[string]interface{}), path[1:])
		default:
			return next
		}
	} else {
		return nil
	}
}

// Given a key, find the value
func (c RawConfig) find(key string) interface{} {
	var val interface{}
	var exists bool

	val, exists = c[key]
	if exists {
		return val
	}
	return nil
}

// Get can retrieve any value given the key to use
// Get returns an interface. For a specific value use one of the Get____ methods.
func (c RawConfig) Get(key string) interface{} {
	path := strings.Split(key, keyDelim)

	val := c.find(strings.ToLower(key))

	if val == nil {
		source := c.find(path[0])
		if source == nil {
			return nil
		}

		if reflect.TypeOf(source).Kind() == reflect.Map {
			val = c.searchMap(cast.ToStringMap(source), path[1:])
		}
	}

	switch val.(type) {
	case bool:
		return cast.ToBool(val)
	case string:
		return cast.ToString(val)
	case int64, int32, int16, int8, int:
		return cast.ToInt(val)
	case float64, float32:
		return cast.ToFloat64(val)
	case time.Time:
		return cast.ToTime(val)
	case time.Duration:
		return cast.ToDuration(val)
	case []string:
		return val
	}
	return val
}

// Returns the value associated with the key as a string
func (c RawConfig) GetString(key string) string {
	return cast.ToString(c.Get(key))
}

// Returns the value associated with the key asa boolean
func (c RawConfig) GetBool(key string) bool {
	return cast.ToBool(c.Get(key))
}

// Returns the value associated with the key as an integer
func (c RawConfig) GetInt(key string) int {
	return cast.ToInt(c.Get(key))
}

// Returns the value associated with the key as a float64
func (c RawConfig) GetFloat64(key string) float64 {
	return cast.ToFloat64(c.Get(key))
}

// Returns the value associated with the key as time
func (c RawConfig) GetTime(key string) time.Time {
	return cast.ToTime(c.Get(key))
}

// Returns the value associated with the key as a duration
func (c RawConfig) GetDuration(key string) time.Duration {
	return cast.ToDuration(c.Get(key))
}

// Returns the value associated with the key as a slice of strings
func (c RawConfig) GetStringSlice(key string) []string {
	return cast.ToStringSlice(c.Get(key))
}

// Returns the value associated with the key as a map of interfaces
func (c RawConfig) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(c.Get(key))
}

// Returns the value associated with the key as a map of strings
func (c RawConfig) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(c.Get(key))
}

// Returns the value associated with the key as a map to a slice of strings.
func (c RawConfig) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(c.Get(key))
}

type SymptomConfig struct {
	Name   string    // Used to lookup the syptom in the plugin registry
	Config RawConfig // Provided to the plugin on load to validate
}

type Config struct {
	Port     string
	Name     string
	Symptoms []SymptomConfig //`yaml:"symptoms,inline`
}

type ConfigLoader struct{}

func (cl *ConfigLoader) Load(data []byte) *Config {
	c := &Config{}

	err := yaml.Unmarshal(data, &c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return c
}
