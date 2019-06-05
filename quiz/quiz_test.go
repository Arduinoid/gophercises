package main

import "testing"

type testStrings struct {
	input    string
	expected string
}

var tests = []testStrings{
	testStrings{" TesTIng  ", "testing"},
	testStrings{"ANOTHER TEST", "another test"},
	testStrings{"One Test With a NUmber 2 in iT   ", "one test with a number 2 in it"},
}

func TestCleanString(t *testing.T) {
	var s string

	for _, test := range tests {
		s = cleanString(test.input)
		if s != test.expected {
			t.Error(
				"For: ", test.input,
				"expected: ", test.expected,
				"got: ", s,
			)
		}
	}
}
