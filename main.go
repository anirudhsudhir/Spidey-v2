package main

import (
	"log"
	"os"

	"github.com/anirudhsudhir/Spidey-v2/crawl"
)

func main() {
	errorLogger := log.New(os.Stdout, "LOG", log.Ldate|log.Ltime|log.Lshortfile)
	parser := Parser{
		ErrorLogger: errorLogger,
	}

	parser.parseArguments()
	parser.parseSeeds()
	// log.Printf("Parser log\n")
	// log.Printf("%+v", parser)

	crawler := crawl.Crawler{
		Seeds:          parser.Seeds,
		TotalCrawlTime: parser.TotalCrawlTime,
		MaxRequestTime: parser.MaxRequestTime,
		ErrorLogger:    errorLogger,
	}

	crawlResults := crawler.StartCrawl()
	log.Printf("Crawler log\n")
	log.Printf("%+v", crawlResults)
}
