package config

import (
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
