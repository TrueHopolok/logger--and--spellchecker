package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	spch "github.com/TrueHopolok/spellchecker/spellchecker"

	plog "github.com/TrueHopolok/plog"
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
	logger.Line()
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
		logger.Fatal("%s", err.Error())
	}
	logger.Info("file '%s' was opened for reading", dictFileName)

	// Reading dictionary
	dictionary := make(map[string]void)
	sdict := bufio.NewScanner(fdict)
	sdict.Split(bufio.ScanWords)
	for sdict.Scan() {
		dictionary[sdict.Text()] = vmem
		logger.Debug("%s", "read: "+sdict.Text())
	}
	logger.Info("all records from file '%s' have been read", dictFileName)

	// Closing dictionary
	err = fdict.Close()
	if err != nil {
		logger.Error("%s", err.Error())
	} else {
		logger.Info("file '%s' was closed", dictFileName)
	}
	logger.Info("program's startup complete, waiting for user input")

	// Getting user input
	var query string
	fmt.Printf("Enter the word to spell check:\n>> ")
	fmt.Scan(&query)
	query = strings.ToLower(query)
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
		curScore := spch.FindScore(word, query)
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
