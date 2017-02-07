package dog

import "testing"

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
