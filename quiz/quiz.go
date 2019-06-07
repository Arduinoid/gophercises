// Quiz game that takes a csv file of question/answers and keeps track of right answers
package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Problem struct for holding one line from the csv file
type Problem struct {
	Question string
	Answer   string
}

// Problems is a slice of Problem type
type Problems []Problem

// ProblemSet just embeds the Problems slice type to be able to hang methods off of
type ProblemSet struct {
	Problems
	correctCount, timeLimit int
	random                  bool
}

// Quiz will be used to provide dependencies to run through propblem set
type Quiz struct {
	ProblemSet
	reader  MyReader
	printer MyPrinter
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

func main() {

	// Run the problem set and show the results when finished or time limit occurs
	q.Run(bufio.NewReader(os.Stdin))
}

// result takes the number of correct answers and the total number of questions and outputs a message to the console
func (q *Quiz) result() {
	q.printer.Printf("You answered %d out of %d correct", q.ProblemSet.correctCount, len(q.ProblemSet.Problems))
	os.Exit(0)
}

// cleanString removes trailing and leading whitespace
func cleanString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// randomizeProblems shuffles a slice of slices contianing question and answer strings
func (ps *ProblemSet) randomizeProblems() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(ps.Problems), func(i int, j int) {
		ps.Problems[i], ps.Problems[j] = ps.Problems[j], ps.Problems[i]
	})
}

// evaluateAnswer takes a problem type and will check a given answer and if it matches will increment a count pointer
func (p Problem) evaluateAnswer(count *int, userAnswer string) {
	if strings.Compare(cleanString(p.Answer), cleanString(userAnswer)) == 0 {
		*count++
		print.Printf("Correct!\n\n")
	} else {
		print.Printf("Wrong :(\n-- correct answer: %s \n\n", p.Answer)
	}
}

// getAnswerFromUser will ask the user the question from the problem and return the users answer
func (p Problem) getAnswerFromUser(r MyReader) string {
	print.Println("what is the answer to: " + p.Question + " ?")
	answer, _ := r.ReadString('\n')
	return answer
}

// addProblem appends a problem to the problem set
func (ps *ProblemSet) addProblem(p Problem) {
	ps.Problems = append(ps.Problems, p)
}

// Run will begin asking the user for answers the the given problem set
func (q *Quiz) Run(reader MyReader) {
	if q.ProblemSet.random {
		q.ProblemSet.randomizeProblems()
	}

	q.printer.Println("--- Quiz ---")
	q.printer.Println("------------")
	q.printer.Println("Press Enter to begin")
	reader.ReadString('\n')

	done := make(chan bool)
	q.startCountDown(done)
	for _, problem := range q.ProblemSet.Problems {
		answer := problem.getAnswerFromUser(reader)
		problem.evaluateAnswer(&q.ProblemSet.correctCount, answer)
		go func() {
			if <-done {
				q.result()
			}
		}()
	}
	q.result()
}

// newProblem creates a problem type from a line or string slice
func newProblem(l []string) Problem {
	return Problem{Question: l[0], Answer: l[1]}
}

func (ps *ProblemSet) getProblemsFromCSV(n string) error {
	// open file
	f, err := os.Open(n)
	if err != nil {
		return err
	}
	defer f.Close()

	// read in file
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return err
	}

	// build out the problem set from the given csv output
	for _, line := range lines {
		p := newProblem(line)
		ps.addProblem(p)
	}

	return nil
}

func (q *Quiz) startCountDown(c chan bool) {
	// setup and start the count down for the quiz
	countDown := time.NewTimer(time.Duration(ps.timeLimit) * time.Second)
	go func() {
		<-countDown.C
		print.Println("Times up!")
		c <- true
		q.result()
	}()
}
