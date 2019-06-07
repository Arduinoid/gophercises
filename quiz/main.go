// Quiz game that takes a csv file of question/answers and keeps track of right answers
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var q Quiz

func init() {
	// catpure CLI flags to configure quiz
	timeLimit := flag.Int("time-limit", 30, "sets the time limit of the quiz")
	randomize := flag.Bool("randomize", false, "Randomize the order of questions")
	filename := flag.String("path", "./problems.csv", "location of question and answer csv file")
	flag.Parse()

	// basic config for the quiz problem set
	q.ProblemSet = ProblemSet{
		timeLimit: *timeLimit,
		random:    *randomize,
	}
	// setup input and output dependencies
	q.printer = new(printer)
	q.reader = bufio.NewReader(os.Stdin)

	// get the and populate the problem set
	err := q.ProblemSet.getProblemsFromCSV(*filename)
	if err != nil {
		os.Exit(1)
	}

}

// main entry point to run program
func main() {
	// Run the problem set and show the results when finished or time limit occurs
	q.Run()
}

// ---------------------------
// UTILITY METHODS AND TYPES
// ---------------------------

// cleanString removes trailing and leading whitespace
func cleanString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// MyReader will allow for mocking ReadString method
type MyReader interface {
	ReadString(delim byte) (string, error)
}

// MyPrinter allows the fmt printer to be mocked
type MyPrinter interface {
	Printf(format string, a ...interface{}) (n int, err error)
	Println(a ...interface{}) (n int, err error)
}

type printer struct{}

func (p printer) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf(format, a...)
}

func (p printer) Println(a ...interface{}) (n int, err error) {
	return fmt.Println(a...)
}
