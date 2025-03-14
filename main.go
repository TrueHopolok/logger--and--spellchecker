package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	plog "spellchecker/prettylogger"
	"syscall"
)

const (
	logFileName  = "p.log"
	dictFileName = "dictionary.txt"
)

type void struct{}

var vmem void

func main() {
	// Init logger
	flog, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer flog.Close()
	logger, _ := plog.NewLogger(plog.LevelInfo, flog, plog.RequireTimestamp|plog.RequireLevel, false)
	logger.Info("program started / logger initializated")

	// Interruption alert
	chanel_interupt := make(chan os.Signal, 1)
	signal.Notify(chanel_interupt, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-chanel_interupt
		logger.Fatal("program was interupted by keyboard")
	}()

	// Opening dictionary
	fdict, err := os.Open(dictFileName)
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info("file '%s' was opened for reading", dictFileName)

	// Reading dictionary
	dictionary := make(map[string]void)
	sdict := bufio.NewScanner(fdict)
	sdict.Split(bufio.ScanWords)
	for sdict.Scan() {
		dictionary[sdict.Text()] = vmem
		logger.Debug("read: " + sdict.Text())
	}
	logger.Info("all records from file '%s' have been read", dictFileName)

	// Closing dictionary
	err = fdict.Close()
	if err != nil {
		logger.Error(err.Error())
	} else {
		logger.Info("file '%s' was closed", dictFileName)
	}
	logger.Info("program's startup complete, waiting for user input")

	// Getting user input
	var query string
	fmt.Printf("Enter the word to spell check:\n>> ")
	fmt.Scan(&query)
	logger.Info("user input is '%s', begining to find the score", query)

	// Counting the score
	bestScore := len(query) + 1
	var simillar []string
	for word := range dictionary {
		if word == query {
			logger.Warn("program finished execution since user input is spelled correctly")
			fmt.Println("Given word is spelled correctly!")
			return
		}
		curScore := FindScore(word, query)
		logger.Debug("word='%s', score='%d'", word, curScore)
		if curScore < bestScore {
			bestScore = curScore
			simillar = nil
			logger.Debug("'best score' was updated and slice was reset")
		}
		if curScore == bestScore {
			simillar = append(simillar, word)
			logger.Debug("word was added to a simillar list")
		}
	}
	logger.Info("scoring was finished, now outputing the results")

	// Outputing the results
	fmt.Printf("Those are words with best simillarity score of %d:\n", bestScore)
	for _, word := range simillar {
		fmt.Println(">>", word)
	}
	logger.Info("results outputed, program was finished")
}

func FindScore(s1, s2 string) int {
	if s1 == s2 {
		return 0
	}
	size1, size2 := len(s1)+1, len(s2)+1
	if size1 == 1 {
		return size2 - 1
	}
	if size2 == 1 {
		return size1 - 1
	}
	table := make([][]int, size1)
	for i := range size1 {
		table[i] = make([]int, size2)
		table[i][0] = i
	}
	for j := range size2 {
		table[0][j] = j
	}
	for i := 1; i < size1; i++ {
		for j := 1; j < size2; j++ {
			table[i][j] = min(table[i][j-1], table[i-1][j]) + 1
			if table[i][j] > table[i-1][j-1] {
				table[i][j] = table[i-1][j-1]
				if s1[i-1] != s2[j-1] {
					table[i][j]++
				}
			}
		}
	}
	return table[size1-1][size2-1]
}
