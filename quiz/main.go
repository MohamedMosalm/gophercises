package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const problemsFileName = "problems.csv"

var (
	correctAnswers,
	totalQuestions int
)

func main() {

	var (
		flagTimer            = flag.Duration("t", time.Second*30, "quiz time in seconds")
		flagProblemsFilename = flag.String("f", problemsFileName, "path to the problems CSV files")
		flagShuffle          = flag.Bool("s", false, "shuffle questions")
	)
	flag.Parse()

	fmt.Printf("press enter to start the quiz\nquiz duration : %v", *flagTimer)
	fmt.Scanln()

	f, err := os.Open(*flagProblemsFilename)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	r := csv.NewReader(f)

	records, err := r.ReadAll()

	if err != nil {
		panic(err)
	}

	totalQuestions = len(records)

	if *flagShuffle {
		rand.Shuffle(len(records), func(i, j int) { records[i], records[j] = records[j], records[i] })
	}

	defer func() {
		fmt.Printf("Result %d/%d\n", correctAnswers, totalQuestions)
	}()

	quizDone := quiz(records)
	quizTimer := time.Tick(*flagTimer)

	select {
	case <-quizDone:
	case <-quizTimer:
	}

}

func quiz(records [][]string) chan bool {
	done := make(chan bool)

	go func() {
		for idx, record := range records {
			question, correctAnswer := record[0], record[1]
			fmt.Printf("%d- %s ?\n", idx+1, question)

			var answer string
			if _, err := fmt.Scan(&answer); err != nil {
				panic(err)
			}
			if answer == correctAnswer {
				correctAnswers++
			}
		}
		done <- true
	}()

	return done
}
