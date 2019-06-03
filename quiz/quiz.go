// Quiz game that takes a csv file of question/answers and keeps track of right answers
package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

// Problem struct for holding one line from the csv file
type Problem struct {
	Question string
	Answer   string
}

func main() {

	filename := "./problems.csv"
	answeredCorrectly := 0
	// numberOfQuestions := 0

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

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("--- Quiz ---")
	fmt.Println("------------")
	fmt.Println("Press Enter to begin")
	reader.ReadString('\n')

	countDown := time.NewTimer(30 * time.Second)
	go func() {
		<-countDown.C
		fmt.Println("Times up!")
		result(answeredCorrectly, len(lines))
		os.Exit(0)
	}()

	for _, line := range lines {
		data := Problem{
			Question: line[0],
			Answer:   line[1],
		}
		fmt.Println("what is the answer to: " + data.Question + " ?")
		answer, _ := reader.ReadString('\n')
		answer = strings.Trim(answer, "\n\r\t")

		if strings.Compare(data.Answer, answer) == 0 {
			answeredCorrectly++
			fmt.Printf("Correct!\n\n")
		} else {
			fmt.Printf("Wrong :(\n-- correct answer: %s \n\n", data.Answer)
		}
	}

	result(answeredCorrectly, len(lines))
}

func result(correct int, questions int) {
	fmt.Printf("You answered %d out of %d correct", correct, questions)
}
