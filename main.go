package main

import (
	"log"
	"os"
	"time"

	"github.com/anirudhsudhir/Spidey-v2/crawl"
)

func main() {
	startTime := time.Now()

	errorLogger := log.New(os.Stdout, "LOG", log.Ldate|log.Ltime|log.Lshortfile)
	parser := Parser{
		ErrorLogger: errorLogger,
	}

	parser.parseArguments()
	parser.parseSeeds()

	crawler := crawl.CrawlerConfig{
		Seeds:          parser.Seeds,
		CrawlTime:      parser.CrawlTime,
		RequestDelay:   parser.RequestDelay,
		WorkerCount:    parser.WorkerCount,
		CrawlStartTime: startTime,
		ErrorLogger:    errorLogger,
	}

	crawlResults := crawler.StartCrawl()

	parser.WriteCrawlData(crawlResults.CrawledLinks)
	log.Printf("Total Crawls: %d, Successful Crawls: %d, Failed Crawls: %d\n", crawlResults.TotalCrawls, crawlResults.SuccessfulCrawls, crawlResults.FailedCrawls)
	log.Printf("Start time: %q, End time: %q, Duration: %q\n", startTime, time.Now(), time.Since(startTime))
}
