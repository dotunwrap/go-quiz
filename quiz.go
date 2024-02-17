package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

func main() {
	csvFileName := flag.String("csv", "problems.csv", "A CSV file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "The time limit for the quiz in seconds")
	shuffleProblemsFlag := flag.Bool("shuffle", false, "Shuffle the questions")
	flag.Parse()

	file, err := os.Open(*csvFileName)

	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %v\n", *csvFileName))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()

	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}

	if *shuffleProblemsFlag {
		shuffleProblems(lines)
	}

	problems := parseLines(lines)
	correct := 0
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

problemLoop:
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, problem.question)
		answerChannel := make(chan string)

		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerChannel <- answer
		}()

		select {
		case <-timer.C:
			break problemLoop
		case answer := <-answerChannel:
			if answer == problem.answer {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	problems := make([]problem, len(lines))
	for i, line := range lines {
		problems[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return problems
}

func shuffleProblems(lines [][]string) {
	rand.Shuffle(len(lines), func(i, j int) {
		lines[i], lines[j] = lines[j], lines[i]
	})
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
