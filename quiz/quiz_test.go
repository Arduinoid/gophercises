package main

import (
	"io"
	"testing"
)

type testStrings struct {
	input    string
	expected string
}

var mockPrinter interface {
	Println(s ...string)
	Printf(s string, i ...interface{})
}

type testPrinter struct {
	fstring string
	values  []interface{}
	called  bool
}

func (p *testPrinter) Println(a ...interface{}) (int, error) {
	p.values = append(p.values, a...)
	p.called = true
	return 0, nil
}

func (p *testPrinter) Printf(s string, i ...interface{}) (int, error) {
	p.fstring = s
	p.values = append(p.values, i...)
	p.called = true
	return 0, nil
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
	q.printer = new(testPrinter)
	q.ProblemSet = ProblemSet{correctCount: 0}

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
		q.evaluateAnswer(test.Problem, test.Input)
		if q.correctCount != tally+1 {
			t.Errorf("correct answer count got %v expected count = %v", q.ProblemSet.correctCount, tally)
		}
	}

	// reset count for next test
	q.ProblemSet.correctCount = 0
	for _, test := range wrongAnswer {
		q.evaluateAnswer(test.Problem, test.Input)
		if q.correctCount != 0 {
			t.Errorf("for wrong answer expected 0 correct answers, but got %v correct", q.ProblemSet.correctCount)
		}
	}
}

type MockReader struct {
	io.Reader
	testMessage string
}

func (mr *MockReader) ReadString(delim byte) (string, error) {
	return mr.testMessage, nil
}

func TestGetAnswerFromUser(t *testing.T) {

	p := Problem{
		Question: "1",
		Answer:   "One",
	}
	a := make(chan string)
	msg := "Good to Go!\n"
	q.reader = &MockReader{testMessage: msg}
	go q.getAnswerFromUser(p, a)
	got := <-a

	if msg != got {
		t.Errorf("got '%s' expected '%s'", got, msg)
	}
}
