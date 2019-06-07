package main

import (
	"os"
	"time"
)

// Quiz will be used to provide dependencies to propblem set
type Quiz struct {
	ProblemSet
	reader  MyReader
	printer MyPrinter
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
	answer := make(chan string)
	go q.startCountDown(done)

	defer q.result()

Loop:
	for _, problem := range q.ProblemSet.Problems {
		go problem.getAnswerFromUser(reader, answer)

		select {
		case <-done:
			break Loop
		case a := <-answer:
			problem.evaluateAnswer(&q.ProblemSet.correctCount, a)
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
	countDown := time.NewTimer(time.Duration(ps.timeLimit) * time.Second)

	<-countDown.C
	print.Println("Times up!")
	c <- true
}
