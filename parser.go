package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Parser struct {
	Seeds        []string
	CrawlTime    time.Duration
	RequestDelay time.Duration
	WorkerCount  int
	ErrorLogger  *log.Logger
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
			match, _ = strings.CutPrefix(match, "\"")
			match, _ = strings.CutSuffix(match, "\"")
			p.Seeds = append(p.Seeds, match)
		}
	}
	if err = scanner.Err(); err != nil {
		p.ErrorLogger.Fatalf("Error reading seeds.txt: %q\n", err)
	}
	return
}

func (p *Parser) parseArguments() {
	if len(os.Args) < 4 {
		p.ErrorLogger.Fatalf("Invalid number of arguments\n")
	}

	crawlTime, err := strconv.Atoi(os.Args[1])
	if err != nil {
		p.ErrorLogger.Fatalf("Error while reading arguments: %q\n", err)
	}
	requestDelay, err := strconv.Atoi(os.Args[2])
	if err != nil {
		p.ErrorLogger.Fatalf("Error while reading arguments: %q\n", err)
	}
	workerCount, err := strconv.Atoi(os.Args[3])
	if err != nil {
		p.ErrorLogger.Fatalf("Error while reading arguments: %q\n", err)
	}
	p.CrawlTime = time.Duration(crawlTime) * time.Second
	p.RequestDelay = time.Duration(requestDelay) * time.Second
	p.WorkerCount = workerCount
}

func (p *Parser) WriteCrawlData(crawlData map[string]string) {
	file, err := os.Create("crawl_data.csv")
	if err != nil {
		p.ErrorLogger.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	for url, status := range crawlData {
		err := writer.Write([]string{url, status})
		if err != nil {
			p.ErrorLogger.Fatal(err)
		}
	}
	writer.Flush()
}
