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

// MyReader will allow for mocking ReadString method
type MyReader interface {
	ReadString(delim byte) (string, error)
}

func main() {

	var seconds int
	var filename string
	var randomize bool

	answeredCorrectly := 0

	flag.IntVar(&seconds, "time-limit", 30, "sets the time limit of the quiz")
	flag.StringVar(&filename, "path", "./problems.csv", "location of question and answer csv file")
	flag.BoolVar(&randomize, "randomize", false, "Randomize the order of questions")
	flag.Parse()

	// open file
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// read in file
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		panic(err)
	}

	// if randomize flag is set then shuffle the problems
	if randomize {
		lines = randomizeProblems(lines)
	}

	// initialize the prompt of the quiz and wait for user input to begin
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("--- Quiz ---")
	fmt.Println("------------")
	fmt.Println("Press Enter to begin")
	inputReader.ReadString('\n')

	// setup and start the count down for the quiz
	countDown := time.NewTimer(time.Duration(seconds) * time.Second)
	go func() {
		<-countDown.C
		fmt.Println("Times up!")
		result(answeredCorrectly, len(lines))
		os.Exit(0)
	}()

	// read through the various sets of questions and answers and check for correctness
	for _, line := range lines {
		p := Problem{
			Question: line[0],
			Answer:   line[1],
		}

		answer := p.getAnswerFromUser(inputReader)

		p.evaluateAnswer(&answeredCorrectly, answer)
	}

	result(answeredCorrectly, len(lines))
}

// result takes the number of correct answers and the total number of questions and outputs a message to the console
func result(correct int, questions int) {
	fmt.Printf("You answered %d out of %d correct", correct, questions)
}

// cleanString removes trailing and leading whitespace
func cleanString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// randomizeProblems shuffles a slice of slices contianing question and answer strings
func randomizeProblems(p [][]string) [][]string {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(p), func(i int, j int) {
		p[i], p[j] = p[j], p[i]
	})
	return p
}

// evaluateAnswer takes a problem type and will check a given answer and if it matches will increment a count pointer
func (p Problem) evaluateAnswer(count *int, userAnswer string) {
	if strings.Compare(cleanString(p.Answer), cleanString(userAnswer)) == 0 {
		*count++
		fmt.Printf("Correct!\n\n")
	} else {
		fmt.Printf("Wrong :(\n-- correct answer: %s \n\n", p.Answer)
	}
}

func (p Problem) getAnswerFromUser(r MyReader) string {
	fmt.Println("what is the answer to: " + p.Question + " ?")
	answer, _ := r.ReadString('\n')
	return answer
}
