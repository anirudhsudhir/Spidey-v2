package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Parser struct {
	Seeds          []string
	TotalCrawlTime time.Duration
	MaxRequestTime time.Duration
	ErrorLogger    *log.Logger
}

func (p *Parser) parseSeeds() {
	seedFile, err := os.Open("seeds.txt")
	if err != nil {
		p.ErrorLogger.Fatalf("Error while opening seeds.txt: %q\n", err)
	}
	defer seedFile.Close()

	scanner := bufio.NewScanner(seedFile)
	for scanner.Scan() {
		re := regexp.MustCompile(`("http)(.*?)(")`)
		match := re.FindString(scanner.Text())
		if match != "" {
			p.Seeds = append(p.Seeds, match)
		}
	}
	if err = scanner.Err(); err != nil {
		p.ErrorLogger.Fatalf("Error reading seeds.txt: %q\n", err)
	}
	return
}

func (p *Parser) parseArguments() {
	if len(os.Args) < 3 {
		p.ErrorLogger.Fatalf("Invalid number of arguments\n")
	}

	crawlTime, err := strconv.Atoi(os.Args[1])
	if err != nil {
		p.ErrorLogger.Fatalf("Error while reading arguments: %q\n", err)
	}
	requestTime, err := strconv.Atoi(os.Args[2])
	if err != nil {
		p.ErrorLogger.Fatalf("Error while reading arguments: %q\n", err)
	}

	p.TotalCrawlTime = time.Duration(crawlTime) * time.Second
	p.MaxRequestTime = time.Duration(requestTime) * time.Millisecond
}
