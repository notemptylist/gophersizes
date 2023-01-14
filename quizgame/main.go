package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// A struct representing the quiz.
type Quiz struct {
	filename  string
	timelimit int
	questions map[string]string
	score     int
}

// parseFile parses the quiz CSV file and stores the questions
// in the `questions` member.
func (q *Quiz) parseFile() error {
	file, err := os.Open(q.filename)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := csv.NewReader(file)

	questions, err := reader.ReadAll()
	for _, qs := range questions {
		q.questions[qs[0]] = qs[1]
	}
	return err
}

// runQuiz runs the quiz and signals via the channel when the list of
// questions is exhaused.
func (q *Quiz) runQuiz(doneCh chan bool) {
	counter := 1
	reader := bufio.NewReader(os.Stdin)
	for question, answer := range q.questions {
		fmt.Printf("%d. %s", counter, question)
		entry, err := reader.ReadString('\n')
		entry = strings.TrimSuffix(entry, "\r\n")
		if err != nil {
			log.Fatalf("Error scanning input")
		}
		if entry == answer {
			q.score++
		}
		counter++
	}
	doneCh <- true
	close(doneCh)
}

// NewQuiz returns a new Quiz.
func NewQuiz(fname *string, limit *int) *Quiz {
	quiz := &Quiz{
		filename:  *fname,
		timelimit: *limit,
		questions: make(map[string]string),
		score:     0,
	}
	return quiz
}

// scoreQuiz scores the quiz and prints out the results
func (q *Quiz) scoreQuiz() {
	percent := float32(q.score) / float32(len(q.questions)) * 100
	fmt.Printf("\nGot %d out of %d : %3.2f%%", q.score, len(q.questions), percent)
}

func main() {
	csv := flag.String("csv", "problems.csv", "a csv file with quiz questions and answers")
	limit := flag.Int("limit", 30, "time limit for each questions in seconds")
	flag.Parse()

	quiz := NewQuiz(csv, limit)
	err := quiz.parseFile()
	if err != nil {
		log.Fatalf(err.Error())
	}
	doneCh := make(chan bool)
	timer := time.NewTimer(time.Second * time.Duration(quiz.timelimit))
	go quiz.runQuiz(doneCh)
	select {
	case <-doneCh:
	case <-timer.C:
		fmt.Println("Timeout reached.")
	}
	quiz.scoreQuiz()
}
