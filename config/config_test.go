package config

import (
	"fmt"
	"github.com/spf13/cast"
	"reflect"
	"testing"
)

func TestFoo(t *testing.T) {
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
	type Foo struct {
		// requiredy, default
		Name string `required:"true" default:""`
		Port string `required:"true" default:"8080"`
	}

	var validate = func(s interface{}) {
		st := reflect.TypeOf(s)
		for i := 0; i < st.NumField(); i++ {
			field := st.Field(i)
			fmt.Printf("field: %v \n", field)
			fmt.Printf("required: %s, default: %s\n", field.Tag.Get("required"), field.Tag.Get("default"))
			fmt.Printf("field name: %v \n", field.Name)
			fmt.Printf("Field value: %v \n", reflect.ValueOf(field))
			if cast.ToBool(field.Tag.Get("required")) == true && field.Tag.Get("default") == "" {
				fmt.Println("Required field and no default")
			}
		}
	}

	validate(Foo{})
}
