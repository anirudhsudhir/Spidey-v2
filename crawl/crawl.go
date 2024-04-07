package crawl

import (
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
	totalCrawls      atomic.Int32
	successfulCrawls atomic.Int32
	failedCrawls     atomic.Int32
)

type CrawlerConfig struct {
	Seeds          []string
	CrawlTime      time.Duration
	RequestDelay   time.Duration
	WorkerCount    int
	CrawlStartTime time.Time
	InfoLogger     *log.Logger
	ErrorLogger    *log.Logger
}

type CrawlResults struct {
	TotalCrawls      int
	SuccessfulCrawls int
	FailedCrawls     int
	CrawledLinks     map[string]string
}

type Crawler struct {
	requestDelay time.Duration
	latestCrawl  map[string]time.Time
	stopCrawl    bool
	lock         sync.RWMutex
	crawledLinks sync.Map
	InfoLogger   *log.Logger
	ErrorLogger  *log.Logger
}

func (c *CrawlerConfig) StartCrawl() CrawlResults {
	jobs := make(chan string, 1000000)
	crawler := Crawler{
		requestDelay: c.RequestDelay,
		latestCrawl:  make(map[string]time.Time),
		InfoLogger:   c.InfoLogger,
		ErrorLogger:  c.ErrorLogger,
	}

	// Initialising channels with seeds
	for _, url := range c.Seeds {
		jobs <- url
	}

	for range c.WorkerCount {
		go crawler.crawlWorker(jobs)
	}

	time.Sleep(c.CrawlTime)

	crawlResults := CrawlResults{
		TotalCrawls:  int(totalCrawls.Load()),
		FailedCrawls: int(failedCrawls.Load()),
		CrawledLinks: make(map[string]string),
	}
	crawlResults.SuccessfulCrawls = crawlResults.TotalCrawls - crawlResults.FailedCrawls

	crawler.crawledLinks.Range(func(key, value any) bool {
		crawlResults.CrawledLinks[key.(string)] = value.(string)
		return true
	})

	return crawlResults
}

func (c *Crawler) crawlWorker(jobs chan string) {
	for url := range jobs {
		c.pingURL(url, jobs)
	}
}

func (c *Crawler) pingURL(URL string, jobs chan<- string) {
	primaryURL, _ := strings.CutPrefix(URL, "https://")
	primaryURL, _ = strings.CutPrefix(primaryURL, "http://")
	primaryURL = strings.Split(primaryURL, "/")[0]
	primaryURLs := strings.Split(primaryURL, ".")
	if len(primaryURLs) < 2 {
		return
	}
	primaryURL = primaryURLs[len(primaryURLs)-2] + "." + primaryURLs[len(primaryURLs)-1]

	c.lock.RLock()
	lastPingTime, found := c.latestCrawl[primaryURL]
	c.lock.RUnlock()

	if found {
		timeElapsed := time.Now().Sub(lastPingTime)
		if timeElapsed < c.requestDelay {
			time.Sleep(timeElapsed)
		}
	}

	resp, err := http.Get(URL)
	if err != nil {
		c.InfoLogger.Println(err)
		failedCrawls.Add(1)
		c.crawledLinks.Store(URL, "crawl failed")
		return
	} else {
		successfulCrawls.Add(1)
		c.crawledLinks.Store(URL, "crawled successfully")
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

			_, ok := c.crawledLinks.Load(match)
			if !ok {
				jobs <- match
			}
		}
	}
}
