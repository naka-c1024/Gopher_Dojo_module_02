package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var cntCorrect uint

func makeWordList() []string {
	fileObj, err := os.Open("/usr/share/dict/words")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
	defer fileObj.Close()

	var wordList []string
	scanner := bufio.NewScanner(fileObj)
	for scanner.Scan() {
		wordList = append(wordList, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
	return wordList
}

func main() {
	wordList := makeWordList()
	rand.Seed(time.Now().UnixNano())
	go func() {
		for {
			randomLine := rand.Intn(235886) + 1
			fmt.Printf("%s\n-> ", wordList[randomLine])
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				if scanner.Text() == wordList[randomLine] {
					cntCorrect++
				}
				break
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				os.Exit(1)
			}
		}
	}()

	select {
	case <-time.After(30 * time.Second):
		fmt.Printf("\nTime's up! Score: %d\n", cntCorrect)
	}
}
