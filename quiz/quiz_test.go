package main

import (
	"io"
	"testing"
)

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

	// setup
	count := 0

	type testUserAnswer struct {
		Problem
		Input string
	}
	correctAnswer := []testUserAnswer{
		{Problem{Question: "1", Answer: "One"}, "ONE"},
		{Problem{Question: "2", Answer: "two"}, "Two "},
		{Problem{Question: "3", Answer: "thRee"}, "  THree"},
	}

	wrongAnswer := []testUserAnswer{
		{Problem{Question: "1", Answer: "One"}, "  five"},
		{Problem{Question: "2", Answer: "two"}, "Four "},
		{Problem{Question: "3", Answer: "thRee"}, " nine "},
	}

	// assert
	for tally, test := range correctAnswer {
		test.Problem.evaluateAnswer(&count, test.Input)
		if count != tally+1 {
			t.Errorf("correct answer count got %v expected count = %v", count, tally)
		}
	}

	// reset count for next test
	count = 0
	for _, test := range wrongAnswer {
		test.Problem.evaluateAnswer(&count, test.Input)
		if count != 0 {
			t.Errorf("for wrong answer expected 0 correct answers, but got %v correct", count)
		}
	}
}

type MockReader struct {
	io.Reader
	testMessage string
}

func (mr *MockReader) ReadString(delim byte) (string, error) {
	return "Good to Go!", nil
}

func TestGetAnswerFromUser(t *testing.T) {

	p := Problem{
		Question: "1",
		Answer:   "One",
	}

	r := new(MockReader)
	r.testMessage = "Good to Go!"
	answer := p.getAnswerFromUser(r)

	if answer != r.testMessage {
		t.Errorf("got '%s' expected '%s'", answer, r.testMessage)
	}
}
