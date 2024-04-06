package crawl

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Crawler struct {
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

func (c *Crawler) StartCrawl() CrawlResults {
	crawlResults := CrawlResults{}
	for _, url := range c.Seeds {
		c.pingURL(url)
		crawlResults.TotalCrawls++
		// log.Printf("%q - %d\n", url, statusCode)
	}
	return crawlResults
}

func (c *Crawler) pingURL(URL string) {
	resp, err := http.Get(URL)
	if err != nil {
		c.ErrorLogger.Fatal(err)
	}
	c.readLinks(resp)
}

func (c *Crawler) readLinks(resp *http.Response) (urls []string) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.ErrorLogger.Fatal(err)
	}

	re := regexp.MustCompile(`("http)(.*?)(")`)
	match := re.FindString(string(body))
	if match != "" {
		match, _ = strings.CutPrefix(match, "\"")
		match, _ = strings.CutSuffix(match, "\"")
		urls = append(urls, match)
		fmt.Println(match)
	}

	return
}
