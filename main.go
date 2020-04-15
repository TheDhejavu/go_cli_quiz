package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Questions struct {
	Question string
	Answer   int
}

func check(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func main() {
	var que []Questions
	var correct, wrong int = 0, 0
	var timedUp bool = false

	fmt.Println("Simple Shell Quiz")
	fmt.Println("---------------------")
	problemFile := flag.String("f", "/problems.csv", "This is the problem path")
	timeLimit := flag.Int64("t", 30, "This is the time limit for the quiz")

	flag.Parse()
	dir, err := os.Getwd()
	check(err)

	problemPath := filepath.Join(dir, *problemFile)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	file, err := os.Open(problemPath)
	defer file.Close()
	check(err)

	csv := csv.NewReader(file)
	for {
		val, err := csv.Read()
		if err == io.EOF {
			break
		}
		check(err)

		answer, err := strconv.Atoi(val[1])
		question := val[0]
		que = append(que, Questions{
			Question: question,
			Answer:   answer,
		})
	}

	QuizLoop:
		for _, value := range que {
			fmt.Printf("Question %s = ", value.Question)
			in := bufio.NewScanner(os.Stdin)
			in.Scan()
			ans, _ := strconv.Atoi(in.Text())

			select {
			case <-timer.C:
				timedUp = true
				break QuizLoop
			default:
				if ans == value.Answer {
					correct++
					fmt.Println("CORRECT")
				} else {
					fmt.Println("WRONG")
					wrong++
				}
			}
		}

	fmt.Printf("\nYou got %d questions correctly while you got %d wrongly\n out of %d", correct, wrong, len(que))
	if timedUp {
		fmt.Println("\nWhooops, Time up")
	}
}
