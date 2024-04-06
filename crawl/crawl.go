package crawl

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	totalCrawls         atomic.Int32
	successfulCrawls    atomic.Int32
	failedCrawls        atomic.Int32
	requestTimeExceeded atomic.Int32
)

type CrawlerConfig struct {
	Seeds          []string
	TotalCrawlTime time.Duration
	MaxRequestTime time.Duration
	ErrorLogger    *log.Logger
}

type Crawler struct {
	TotalCrawls         int
	SuccessfulCrawls    int
	FailedCrawls        int
	RequestTimeExceeded int
	ErrorLogger         *log.Logger
}

func (c *CrawlerConfig) StartCrawl() Crawler {
	jobs := make(chan string, 100)
	crawler := Crawler{ErrorLogger: c.ErrorLogger}
	numWorkers := 5
	var wg sync.WaitGroup

	// Initialising channels with seeds
	for _, url := range c.Seeds {
		jobs <- url
	}

	for range numWorkers {
		wg.Add(1)
		go crawler.crawlWorker(jobs, &wg)
	}

	for {
		if totalCrawls.Load() == 10 {
			close(jobs)
			break
		}
	}

	wg.Wait()
	crawler.TotalCrawls = int(totalCrawls.Load())
	crawler.SuccessfulCrawls = int(successfulCrawls.Load())
	crawler.FailedCrawls = int(failedCrawls.Load())
	crawler.RequestTimeExceeded = int(requestTimeExceeded.Load())
	return crawler
}

func (c *Crawler) crawlWorker(jobs chan string, wg *sync.WaitGroup) {
	for url := range jobs {
		c.pingURL(url, jobs)
	}
	wg.Done()
}

func (c *Crawler) pingURL(URL string, jobs chan<- string) {
	resp, err := http.Get(URL)
	if err != nil {
		c.ErrorLogger.Fatal(err)
	}
	totalCrawls.Add(1)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.ErrorLogger.Fatal(err)
	}

	// Searching for links from current page
	re := regexp.MustCompile(`("http)(.*?)(")`)
	matches := re.FindAllString(string(body), -1)

	for _, match := range matches {
		if match != "" {
			match, _ = strings.CutPrefix(match, "\"")
			match, _ = strings.CutSuffix(match, "\"")
			fmt.Println(match)
			jobs <- match
		}
	}
}
