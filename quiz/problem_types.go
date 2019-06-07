package main

import (
	"encoding/csv"
	"math/rand"
	"os"
	"time"
)

// Problem struct for holding one line from the csv file
type Problem struct {
	Question string
	Answer   string
}

// ProblemSet just embeds the Problems slice type to be able to hang methods off of
type ProblemSet struct {
	Problems                []Problem
	correctCount, timeLimit int
	random                  bool
}

// randomizeProblems shuffles a slice of slices contianing question and answer strings
func (ps *ProblemSet) randomizeProblems() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(ps.Problems), func(i int, j int) {
		ps.Problems[i], ps.Problems[j] = ps.Problems[j], ps.Problems[i]
	})
}

// addProblem appends a problem to the problem set
func (ps *ProblemSet) addProblem(p Problem) {
	ps.Problems = append(ps.Problems, p)
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
		ps.addProblem(newProblem(line))
	}

	return nil
}
