package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type task struct {
	q string
	a string
}

func taskPuller(fileName string) ([]task, error) {
	//read tasks from .csv file
	//open the file
	if fObj, err := os.Open(fileName); err == nil {
		//create new reader
		csvR := csv.NewReader(fObj)
		//read the file
		if cLines, err := csvR.ReadAll(); err == nil {
			return parseTask(cLines), nil
		} else {
			return nil, fmt.Errorf("error in reading data in csv from %s file; %s", fileName, err.Error())
		}
	} else {
		return nil, fmt.Errorf("error in opening %s file; %s", fileName, err.Error())
	}

}

func parseTask(lines [][]string) []task {
	//go over the lines and parse them

	r := make([]task, len(lines))
	for i := 0; i < len(lines); i++ {
		r[i] = task{q: lines[i][0], a: lines[i][1]}
	}
	return r

}

func main() {
	//input the name of file
	fName := flag.String("f", "quiz.csv", "path of csv file")
	//set duration of the timer
	timer := flag.Int("t", 30, "timer for the quiz")
	flag.Parse()
	//pull task from the file
	tasks, err := taskPuller(*fName)
	//handle the error
	if err != nil {
		exit(fmt.Sprintf("error:%s", err.Error()))
	}

	//create counter for correct answers
	correctAns := 0
	//initialize the timer
	tObj := time.NewTimer(time.Duration(*timer) * time.Second)
	ansCh := make(chan string)
	//loop through the questions, print the question, accept the answers
taskLoop:
	for i, t := range tasks {
		var answer string
		fmt.Printf("Task %d: %s=", i+1, t.q)

		go func() {
			fmt.Scan(&answer)
			ansCh <- answer
		}()
		select {
		case <-tObj.C:
			fmt.Println()
			break taskLoop
		case iAns := <-ansCh:
			if iAns == t.a {
				correctAns++
			}
			if i == len(tasks)-1 {
				close(ansCh)
			}
		}
	}
	//calculate and print out the result
	fmt.Printf("Corect answers is %d from %d\n", correctAns, len(tasks))
	fmt.Printf("Press enter to exit")
	//<-ansCh
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
