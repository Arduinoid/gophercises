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

func TestEvaluateAnswer(t *testing.T) {

	count := 0

	type testUserAnswer struct {
		Problem
		Input string
	}
	tests := []testUserAnswer{
		{Problem{Question: "1", Answer: "One"}, "ONE"},
		{Problem{Question: "2", Answer: "two"}, "Two "},
		{Problem{Question: "3", Answer: "thRee"}, "  THree"},
	}

	for tally, test := range tests {
		test.Problem.evaluateAnswer(&count, test.Input)
		if count != tally+1 {
			t.Errorf("count = %v expected count = %v", count, tally)
		}
	}
}
