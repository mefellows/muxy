package config

import (
	"fmt"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	var data = []byte(`
port: 8080
name: Foo
symptoms:
  - name: hell
    config:
      foo: bar
      baz: bat
      bar:
        - 1
        - 2
        - 3
  - name: fire
    config:
      foo: bar
      baz: bat
`)

	cl := &ConfigLoader{}
	c, _ := cl.Load([]byte(data))

	if len(c.Symptoms) != 2 {
		t.Fatalf("Expected 2 symptoms, but got %d", len(c.Symptoms))
	}
	if c.Symptoms[0].Name != "hell" {
		t.Fatalf("Expected the first symptom to have the name 'hell' but got '%s'", c.Symptoms[0].Name)
	}
	if !reflect.DeepEqual(c.Symptoms[0].Config.GetStringSlice("bar"), []string{"1", "2", "3"}) {
		t.Fatalf("Expected 'bar' property of symptoms[0].Config.bar to equal [1 2 3] but got: %v", c.Symptoms[0].Config.Get("bar"))
	}
}

func TestRawConfigToStruct(t *testing.T) {
	var data = []byte(`
port: 8080
name: Foo
symptoms:
  - name: fire
    config:
      name: billy 
      foo: bar
`)

	var input interface{}
	cl := &ConfigLoader{}
	c, _ := cl.Load([]byte(data))
	input = c.Symptoms[0].Config

	type Foo struct {
		Name string `required:"true" default:"" mapstructure:"name"`
		Port int    `required:"true" default:"8081" mapstructure:"port"`
		Foo  string `required:"true"`
	}

	var f Foo

	err := cl.ApplyConfig(input, &f)
	if err != nil {
		t.Fatalf("Did not expect error: %v", err)
	}
	err = cl.Validate(f, &f)
	fmt.Printf("Result: %v", f)
	if err != nil {
		t.Fatalf("Did not expect error: %v", err)
	}

	expected := Foo{Name: "billy", Port: 8081, Foo: "bar"}
	if !reflect.DeepEqual(f, expected) {
		t.Fatalf("Expected %v but got %v", expected, f)
	}
}

func TestRawConfigToStructWithErrors(t *testing.T) {
	var data = []byte(`
port: 8080
name: Foo
symptoms:
  - name: fire
    config:
      name: bar
      port: 1234
`)

	var input interface{}
	cl := &ConfigLoader{}
	c, _ := cl.Load([]byte(data))
	input = c.Symptoms[0].Config

	type Foo struct {
		Name string `required:"true" default:"" mapstructure:"name"`
		Port int    `required:"true" default:"8081" mapstructure:"port"`
		Foo  string `required:"true"`
	}

	var f Foo = Foo{}

	err := cl.ApplyConfig(input, &f)
	if err != nil {
		t.Fatalf("Did not expect error: %v", err)
	}
	err = cl.Validate(f, &f)
	if err == nil {
		t.Fatalf("Expected error, but did not get one")
	}
}
