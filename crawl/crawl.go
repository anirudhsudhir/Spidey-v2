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

type CrawlResults struct {
	TotalCrawls         int
	SuccessfulCrawls    int
	FailedCrawls        int
	RequestTimeExceeded int
}

type Crawler struct {
	latestCrawl map[string]time.Time
	stopCrawl   bool
	lock        sync.RWMutex
	ErrorLogger *log.Logger
}

func (c *CrawlerConfig) StartCrawl() CrawlResults {
	jobs := make(chan string, 100)
	crawler := Crawler{
		latestCrawl: make(map[string]time.Time),
		ErrorLogger: c.ErrorLogger,
	}
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
			crawler.lock.Lock()
			crawler.stopCrawl = true
			crawler.lock.Unlock()

			time.Sleep(time.Second)
			close(jobs)
			break
		}
	}

	wg.Wait()
	crawlResults := CrawlResults{
		TotalCrawls:         int(totalCrawls.Load()),
		SuccessfulCrawls:    int(successfulCrawls.Load()),
		FailedCrawls:        int(failedCrawls.Load()),
		RequestTimeExceeded: int(requestTimeExceeded.Load()),
	}
	return crawlResults
}

func (c *Crawler) crawlWorker(jobs chan string, wg *sync.WaitGroup) {
	for url := range jobs {
		c.pingURL(url, jobs)
	}
	wg.Done()
}

func (c *Crawler) pingURL(URL string, jobs chan<- string) {
	// fmt.Println("EDITING", URL)
	primaryURL, _ := strings.CutPrefix(URL, "https://")
	primaryURL, _ = strings.CutPrefix(primaryURL, "http://")
	// fmt.Println("EDITING", primaryURL)
	primaryURL = strings.Split(primaryURL, "/")[0]
	// fmt.Println("EDITING", primaryURL)
	primaryURLs := strings.Split(primaryURL, ".")
	fmt.Println("EDITING", primaryURL)
	if len(primaryURLs) < 2 {
		return
	}
	primaryURL = primaryURLs[len(primaryURLs)-2] + "." + primaryURLs[len(primaryURLs)-1]

	c.lock.RLock()
	lastPingTime, found := c.latestCrawl[primaryURL]
	c.lock.RUnlock()

	if found {
		timeElapsed := time.Now().Sub(lastPingTime)
		if timeElapsed < time.Duration(time.Second) {
			time.Sleep(timeElapsed)
			fmt.Println("slept ", timeElapsed, " ", primaryURL)
		}
	}

	resp, err := http.Get(URL)
	if err != nil {
		c.ErrorLogger.Println(err)
	}
	c.lock.Lock()
	c.latestCrawl[primaryURL] = time.Now()
	c.lock.Unlock()

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
			// fmt.Println("formatting", match)

			c.lock.RLock()
			stopCrawl := c.stopCrawl
			c.lock.RUnlock()
			if match != "http://" && match != "https://" && match != "http:" && match != "https:" && !stopCrawl {
				fmt.Println("stopCrawl: ", stopCrawl)
				jobs <- match
			}
		}
	}
}
