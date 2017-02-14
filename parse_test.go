package dog

import (
	"reflect"
	"testing"
)

func TestValidDogfileName(t *testing.T) {
	for i, test := range []struct {
		input  string
		expect bool
	}{
		{"Dogfile.yml", true},
		{"Dogfile.yaml", true},
		{"Dogfile", true},
		{"🐕.yml", true},
		{"Dogfile-foo.yml", true},
		{"dogfile.yml", false},
		{"DogFile.yml:", false},
	} {
		if got, want := validDogfileName(test.input), test.expect; got != want {
			t.Errorf("Test %d (%s): expected %v but was %v", i, test.input, want, got)
		}
	}
}

func TestValidTaskName(t *testing.T) {
	for i, test := range []struct {
		input  string
		expect bool
	}{
		{"foo", true},
		{"foo-bar", true},
		{"01-with-02-numbers-03", true},
		{"-foo", false},
		{"foo-", false},
		{"-", false},
		{"camelCase", false},
		{"snake_case:", false},
		{"Some-Caps", false},
	} {
		if got, want := validTaskName(test.input), test.expect; got != want {
			t.Errorf("Test %d (%s): expected %v but was %v", i, test.input, want, got)
		}
	}
}

func TestDogfileParseYAML(t *testing.T) {
	got, err := Parse([]byte(`
- task: foo
  description: Foo task
  post: bar
  code: echo "foo"

- task: bar
  description: Bar task
  code: echo "bar"
`))
	if err != nil {
		t.Fatalf("Failed parsing Dogfile from YAML: %v", err)
	}

	want := Dogfile{
		Tasks: map[string]*Task{
			"foo": {
				Name:        "foo",
				Description: "Foo task",
				Post:        []string{"bar"},
				Code:        "echo \"foo\"",
			},
			"bar": {
				Name:        "bar",
				Description: "Bar task",
				Code:        "echo \"bar\"",
			},
		},
	}

	if want.Tasks["foo"].Name != got.Tasks["foo"].Name ||
		want.Tasks["foo"].Description != got.Tasks["foo"].Description ||
		want.Tasks["foo"].Post[0] != got.Tasks["foo"].Post[0] ||
		want.Tasks["foo"].Code != got.Tasks["foo"].Code ||
		want.Tasks["bar"].Name != got.Tasks["bar"].Name ||
		want.Tasks["bar"].Description != got.Tasks["bar"].Description ||
		want.Tasks["bar"].Code != got.Tasks["bar"].Code {
		t.Fatalf("Expected %v but was %v", want, got)
	}
}

func TestDogfileParseJSON(t *testing.T) {
	if _, err := Parse([]byte(`[
  {
    "task": "foo",
    "description": "Foo task",
    "post": "bar",
    "code": "echo \"foo\""
  },
  {
    "task": "bar",
    "description": "Bar task",
    "code": "echo \"bar\""
  }
]`)); err != nil {
		t.Errorf("Failed parsing Dogfile from JSON: %s", err)
	}
}

func TestDogfileParseDuplicatedTask(t *testing.T) {
	if _, err := Parse([]byte(`
- task: foo
  description: Foo task
  code: echo "foo"

- task: foo
  description: Foo task
  code: echo "foo"
`)); err == nil {
		t.Errorf("Failed to detect duplicated task name")
	}
}

func TestDogfileParsePreTasksArray(t *testing.T) {
	dogfile, err := Parse([]byte(`
- task: lorem
  description: Foo task
  pre:
    - ipsum
    - dolor
  code: echo "lorem"

- task: ipsum
  code: echo "ipsum"

- task: dolor
  code: echo "dolor"
`))
	if err != nil {
		t.Fatalf("Failed to parse pre tasks array: %v", err)
	}

	got := dogfile.Tasks["lorem"].Pre
	want := []string{"ipsum", "dolor"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected %v but was %v", want, got)
	}
}

func TestDogfileValidatePost(t *testing.T) {
	dogfile := Dogfile{
		Tasks: map[string]*Task{
			"foo": {
				Name:        "foo",
				Description: "Foo is a task that, it will trigger Bar when it finishes",
				Post:        []string{"bar"},
				Code:        "echo foo",
			},
			"bar": {
				Name: "bar",
				Code: "echo bar",
			},
		},
	}
	err := dogfile.Validate()
	if err != nil {
		t.Errorf("Failed validating a Dogfile with a post task: %v", err)
	}
}

func TestDogfileValidatePostError(t *testing.T) {
	dogfile := Dogfile{
		Tasks: map[string]*Task{
			"foo": {
				Name:        "foo",
				Description: "Foo is a task that, it will trigger Bar when it finishes",
				Post:        []string{"bar"},
				Code:        "echo foo",
			},
		},
	}
	err := dogfile.Validate()
	if err == nil {
		t.Errorf("Failed, should have errored validating a Dogfile with an unexistent post task")
	}
}
