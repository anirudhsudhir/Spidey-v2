package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

func main() {
	totalCrawlTime, maxRequestTime := parseArguments()
	seedUrls := parseSeeds()
	// crawler := crawl.Crawler{
	// 	Seeds:          seedUrls,
	// 	TotalCrawlTime: totalCrawlTime,
	// 	MaxRequestTime: maxRequestTime,
	// }
	// crawlStats := crawl.CrawlLinks(seedUrls, totalCrawlTime, maxRequestTime)
	// log.Printf("TotalCrawls: %d, SuccessfulCrawls: %d, FailedCrawls: %d, Request Time Exceeded: %d", crawlStats.TotalCrawls, crawlStats.SuccessfulCrawls, crawlStats.FailedCrawls, crawlStats.RequestTimeExceeded)
	log.Printf("TotalCrawlTime:%d, MaxRequestTime:%d, seedUrls:%v", totalCrawlTime, maxRequestTime, seedUrls)
}

func parseSeeds() (seedUrls []string) {
	seedFile, err := os.Open("seeds.txt")
	if err != nil {
		log.Fatalf("Error while opening seeds.txt: %q\n", err)
	}
	defer seedFile.Close()

	scanner := bufio.NewScanner(seedFile)
	for scanner.Scan() {
		re := regexp.MustCompile(`("http)(.*?)(")`)
		match := re.FindString(scanner.Text())
		if match != "" {
			seedUrls = append(seedUrls, match)
		}
	}
	if err = scanner.Err(); err != nil {
		log.Fatalf("Error reading seeds.txt: %q\n", err)
	}
	return
}

func parseArguments() (time.Duration, time.Duration) {
	if len(os.Args) < 3 {
		log.Fatalf("Invalid number of arguments")
	}

	crawlTime, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Error while reading arguments: %q", err)
	}
	requestTime, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("Error while reading arguments: %q", err)
	}

	totalCrawlTime := time.Duration(crawlTime) * time.Second
	maxRequestTime := time.Duration(requestTime) * time.Millisecond
	return totalCrawlTime, maxRequestTime
}
