package crawl

import (
	"log"
	"time"
)

type Crawler struct {
	Seeds          []string
	TotalCrawlTime time.Duration
	MaxRequestTime time.Duration
	ErrorLogger    log.Logger
}

type CrawlResults struct {
	TotalCrawls         int
	SuccessfulCrawls    int
	FailedCrawls        int
	RequestTimeExceeded int
}

func (c *Crawler) StartCrawl() CrawlResults {
	return CrawlResults{}
}
