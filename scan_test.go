package main

import "testing"

type testData struct {
	input    []string
	expected string
}

func Test(t *testing.T) {

	tests := []testData{
		{input: []string{"127.0.0.1", "80"}, expected: "[127.0.0.1 80]"},
	}

	for _, test := range tests {
		actual := scan(test.input)
		if actual != test.expected {
			t.Errorf("Expected: %s; Actual: %s", test.expected, actual)
		}
	}
}
