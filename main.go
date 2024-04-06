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
	// log.Printf("Parser log\n")
	// log.Printf("%+v", parser)

	crawler := crawl.CrawlerConfig{
		Seeds:          parser.Seeds,
		TotalCrawlTime: parser.TotalCrawlTime,
		MaxRequestTime: parser.MaxRequestTime,
		ErrorLogger:    errorLogger,
	}

	crawlResults := crawler.StartCrawl()
	log.Printf("Crawler log\n")
	log.Printf("%+v", crawlResults)
	log.Printf("Start time: %q, End time: %q, Duration: %q\n", startTime, time.Now(), time.Since(startTime))
}
