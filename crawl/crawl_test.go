package crawl_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/anirudhsudhir/Spidey-v2/crawl"
)

type testConfig struct {
	urlSetSize    int
	seedServers   []*httptest.Server
	testUrls      []string
	allServers    []*httptest.Server
	secondaryUrls string
	tempServer    *httptest.Server
}

func TestPingWebsites(t *testing.T) {
	testconfig := testConfig{
		urlSetSize: 5,
	}
	testconfig.initTestServers()

	t.Run("all links crawled", func(t *testing.T) {
		crawler := crawl.Crawler{
			Seeds:          testconfig.testUrls,
			TotalCrawlTime: time.Duration(5 * time.Second),
			MaxRequestTime: time.Duration(1 * time.Millisecond),
			ErrorLogger:    log.New(os.Stdout, "LOG", log.Ldate|log.Ltime|log.Lshortfile),
		}

		got := crawler.StartCrawl()
		want := testconfig.urlSetSize + testconfig.urlSetSize*testconfig.urlSetSize
		testconfig.stopTestServers()

		if got.TotalCrawls != want {
			t.Errorf("crawled %d links, want %d links", got, want)
		}
	})
}

func (t *testConfig) initTestServers() {
	for i := 0; i < t.urlSetSize; i++ {
		t.secondaryUrls = "random text"
		for j := 0; j < t.urlSetSize; j++ {
			t.tempServer = createServer(time.Duration(50)*time.Millisecond, "random text")
			t.secondaryUrls += "\"" + t.tempServer.URL + "\""
			t.allServers = append(t.allServers, t.tempServer)
		}
		t.seedServers = append(t.seedServers, createServer(50*time.Millisecond, t.secondaryUrls))
		t.allServers = append(t.allServers, t.seedServers[i])
		t.testUrls = append(t.testUrls, t.seedServers[i].URL)
	}
}

func (t *testConfig) stopTestServers() {
	for i := 0; i < len(t.allServers); i++ {
		t.allServers[i].Close()
	}
}

func createServer(delay time.Duration, message string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(message))
	}))
}
