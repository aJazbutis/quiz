package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func normaliseString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

type Problem struct {
	Question string
	Answer   string
}

type Result struct {
	Total   int
	Correct int
}

func shuffleProblems(s []Problem) []Problem {
	//autoseeded nowdays
	for i := len(s) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		s[j], s[i] = s[i], s[j]
	}
	return s
}

func parseCsv(path *string, shuffle *bool) ([]Problem, error) {
	ret := []Problem{}
	if file, err := os.Open(*path); err == nil {
		defer file.Close()
		r := csv.NewReader(file)
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			ret = append(ret, Problem{Question: record[0], Answer: record[1]})
		}
		if *shuffle {
			return shuffleProblems(ret), nil
		}
		return ret, nil
	} else {
		return nil, err
	}
}

func cntrlL() {
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func play(p Problem, c chan bool) {
	var s string
	fmt.Println(p.Question)
	fmt.Scanln(&s)
	if normaliseString(s) == normaliseString(p.Answer) {
		c <- true
	} else {
		c <- false
	}
}

func main() {
	path := flag.String("csv", "../problems.csv", "path to .csv file, the csv format in 'question, answer'")
	t := flag.Int("t", 30, "time limit")
	shuffle := flag.Bool("shuffle", false, "shuffle the order of questions")
	flag.Parse()
	if *t < 0 {
		log.Fatal("Negative time value")
	}
	if problems, err := parseCsv(path, shuffle); err == nil {
		var s string
		answer := make(chan bool)
		var result Result
		fmt.Println("Start?")
		fmt.Scanln(&s)
		if normaliseString(s) == "no" {
			return
		}
	game:
		for _, problem := range problems {
			result.Total++
			cntrlL()
			go play(problem, answer)
		round:
			for {
				select {
				case correct := <-answer:
					if correct {
						result.Correct++
					}
					break round
				case <-time.After(time.Second * time.Duration(*t)):
					fmt.Println("Time ran out.")
					break game
				}
			}
		}
		fmt.Printf("%v correct out of %v. Awesome.\n", result.Correct, result.Total)
	} else {
		log.Fatal(err)
	}
}
