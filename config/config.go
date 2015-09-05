package config

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
	"strconv"
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

func (cl *ConfigLoader) Load(data []byte) (*Config, error) {
	c := &Config{}

	err := yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
func (cl *ConfigLoader) LoadFromFile(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return cl.Load(data)
}
func (cl *ConfigLoader) ApplyConfig(config interface{}, i interface{}) error {
	return mapstructure.Decode(config, i)
}

// Give me a struct with field tags and i'll validate you, set defaults, etc.
func (cl *ConfigLoader) Validate(iface interface{}) error {
	iValue := reflect.ValueOf(iface).Elem().Interface()
	st := reflect.TypeOf(iValue)
	ps := reflect.ValueOf(iValue)

	// Loop all fields, set their default values if they have them and are empty
	// Fail if mandatory fields are not set and have no value

	for i := 0; i < ps.NumField(); i++ {
		f := st.Field(i)
		field := ps.FieldByName(f.Name)

		if cast.ToBool(f.Tag.Get("required")) == true && isZero(field) {
			defaultVal := f.Tag.Get("default")
			if defaultVal == "" {
				return errors.New(fmt.Sprintf("Mandatory field '%s' has not been set, and has no provided default", f.Name))
			}

			dataKind := field.Kind()
			field = reflect.ValueOf(iface).Elem().FieldByName(f.Name)
			switch dataKind {
			case reflect.Bool:
				field.SetBool(cast.ToBool(defaultVal))
			case reflect.String:
				field.SetString(defaultVal)
			case reflect.Slice, reflect.Array:
				_type := field.Type().Elem()
				_newArr := strings.Split(f.Tag.Get("default"), ",")

				switch _type {
				case reflect.TypeOf(""):
					field.Set(reflect.ValueOf(_newArr))
				case reflect.TypeOf(1):
					// Convert array to int
					intArray := make([]int, len(_newArr))
					var err error
					for j, val := range _newArr {
						intArray[j], err = strconv.Atoi(val)
						if err != nil {

							return errors.New(fmt.Sprintf("Error creating default array for field '$s': %v\n", f.Name, err))
						}
					}
					field.Set(reflect.ValueOf(intArray))
				default:
					return errors.New(fmt.Sprintf("Unsupported slice default type: %v\n", _type))
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				s, err := strconv.ParseInt(defaultVal, 10, 64)
				if err != nil {
					return err
				}
				field.SetInt(s)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				s, err := strconv.ParseUint(defaultVal, 10, 64)
				if err != nil {
					return err
				}
				field.SetUint(s)
			default:
				return errors.New(fmt.Sprintf("Unsupported field '%s' of type: %s", f.Name, dataKind))
			}

		}
	}
	return nil
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && isZero(v.Index(i))
		}
		return z
	case reflect.Struct:
		z := true
		for i := 0; i < v.NumField(); i++ {
			z = z && isZero(v.Field(i))
		}
		return z
	}
	// Compare other types directly:
	z := reflect.Zero(v.Type())
	return v.Interface() == z.Interface()
}
