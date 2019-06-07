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

var print printer

var ps ProblemSet

func init() {
	var seconds int
	var filename string
	var randomize bool

	flag.IntVar(&seconds, "time-limit", 30, "sets the time limit of the quiz")
	flag.StringVar(&filename, "path", "./problems.csv", "location of question and answer csv file")
	flag.BoolVar(&randomize, "randomize", false, "Randomize the order of questions")
	flag.Parse()

	// get the and populate the problem set
	ps.getProblemsFromCSV(filename)
	ps.random = randomize
	ps.timeLimit = seconds

	q.ProblemSet = ps
	q.printer = print
}

// main entry point to run program
func main() {

	// Run the problem set and show the results when finished or time limit occurs
	q.Run(bufio.NewReader(os.Stdin))
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
