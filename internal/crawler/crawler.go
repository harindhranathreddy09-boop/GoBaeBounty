package crawler

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/common"
)

// Run performs recursive crawling up to depth with robots.txt respect and collects HTML pages and JS files
func Run(ctx context.Context, config *common.Config, discoveryResults *common.DiscoveryResults) (*common.CrawlResults, error) {
	results := &common.CrawlResults{
		Pages:   []common.Page{},
		JSFiles: []common.JSFile{},
		Forms:   []common.Form{},
	}

	robotMgr := NewRobotsManager("GoBaeBounty")

	collyCollector := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(config.CrawlDepth),
		colly.AllowedDomains(discoveryResults.Subdomains...),
	)

	// Respect robots.txt if not ignored
	if !config.IgnoreRobots {
		collyCollector.WithTransport(&robotstxtTransport{robotMgr: robotMgr})
	}

	collyCollector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: config.Workers,
		Delay:       1000 * time.Millisecond, // To respect crawl delays better, consider integrating more intelligent delays
	})

	var mu sync.Mutex

	collyCollector.OnResponse(func(r *colly.Response) {
		// Record pages
		if strings.Contains(r.Headers.Get("Content-Type"), "text/html") {
			page := common.Page{
				URL:        r.Request.URL.String(),
				StatusCode: r.StatusCode,
				Headers:    r.Headers,
				BodySize:   len(r.Body),
				CrawlTime:  time.Since(r.Request.Ctx.GetAny("startTime").(time.Time)).Milliseconds(),
			}
			mu.Lock()
			results.Pages = append(results.Pages, page)
			mu.Unlock()
		}
	})

	collyCollector.OnHTML("script[src]", func(e *colly.HTMLElement) {
		jsURL := e.Request.AbsoluteURL(e.Attr("src"))
		if jsURL == "" {
			return
		}
		mu.Lock()
		results.JSFiles = append(results.JSFiles, common.JSFile{URL: jsURL})
		mu.Unlock()
	})

	collyCollector.OnHTML("form", func(e *colly.HTMLElement) {
		form := common.Form{
			URL:    e.Request.URL.String(),
			Action: e.Attr("action"),
			Method: e.Attr("method"),
		}

		e.ForEach("input", func(i int, el *colly.HTMLElement) {
			form.Inputs = append(form.Inputs, common.FormInput{
				Name:  el.Attr("name"),
				Type:  el.Attr("type"),
				Value: el.Attr("value"),
			})
		})

		mu.Lock()
		results.Forms = append(results.Forms, form)
		mu.Unlock()
	})

	collyCollector.OnRequest(func(r *colly.Request) {
		r.Ctx.Put("startTime", time.Now())
		// Wait crawl delay from robots.txt for domain if not ignored
		if !config.IgnoreRobots {
			delay := robotMgr.CrawlDelay(r.URL.String())
			if delay > 0 {
				time.Sleep(delay)
			}
		}
	})

	// Start crawling from subdomains discovered
	for _, subdomain := range discoveryResults.Subdomains {
		u := "https://" + subdomain
		collyCollector.Visit(u)
	}
	collyCollector.Wait()

	return results, nil
}

type robotstxtTransport struct {
	robotMgr *RobotsManager
}

func (t *robotstxtTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	allowed, err := t.robotMgr.Allowed(req.URL.String())
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, fmt.Errorf("blocked by robots.txt: %s", req.URL.String())
	}
	return http.DefaultTransport.RoundTrip(req)
}
