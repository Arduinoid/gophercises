package main

import (
	"os"
	"strings"
	"time"
)

// Quiz will be used to provide dependencies to propblem set
type Quiz struct {
	ProblemSet
	reader  MyReader
	printer MyPrinter
}

// Run will begin asking the user for answers the the given problem set
func (q *Quiz) Run() {
	if q.ProblemSet.random {
		q.ProblemSet.randomizeProblems()
	}

	q.printer.Println("--- Quiz ---")
	q.printer.Println("------------")
	q.printer.Println("Press Enter to begin")
	q.reader.ReadString('\n')

	done := make(chan bool)
	answer := make(chan string)
	go q.startCountDown(done)

	defer q.result()

Loop:
	for _, problem := range q.ProblemSet.Problems {
		go q.getAnswerFromUser(problem, answer)

		select {
		case <-done:
			break Loop
		case a := <-answer:
			q.evaluateAnswer(problem, a)
		}
	}
}

// result takes the number of correct answers and the total number of questions and outputs a message to the console
func (q *Quiz) result() {
	q.printer.Printf("You answered %d out of %d correct", q.ProblemSet.correctCount, len(q.ProblemSet.Problems))
	os.Exit(0)
}

func (q *Quiz) startCountDown(c chan bool) {
	// setup and start the count down for the quiz
	countDown := time.NewTimer(time.Duration(q.ProblemSet.timeLimit) * time.Second)

	<-countDown.C
	q.printer.Println("Times up!")
	c <- true
}

// evaluateAnswer takes a problem type and will check a given answer and if it matches will increment a count pointer
func (q *Quiz) evaluateAnswer(p Problem, userAnswer string) {
	if strings.Compare(cleanString(p.Answer), cleanString(userAnswer)) == 0 {
		q.correctCount++
		q.printer.Printf("Correct!\n\n")
	} else {
		q.printer.Printf("Wrong :(\n-- correct answer: %s \n\n", p.Answer)
	}
}

// getAnswerFromUser will ask the user the question from the problem and return the users answer
func (q *Quiz) getAnswerFromUser(p Problem, c chan string) {
	q.printer.Println("what is the answer to: " + p.Question + " ?")
	answer, _ := q.reader.ReadString('\n')
	c <- answer
}
