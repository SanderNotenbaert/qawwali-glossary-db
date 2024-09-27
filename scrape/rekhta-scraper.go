package scrape

import (
	"fmt"
	"qawwali-syllabus/translate"
	"slices"
	"strings"

	"github.com/gocolly/colly"
)

//TODO: make concurrent?
// initialize a data structure to keep the scraped data

func Sufinama(domain string, urlList []string, visitedUrls []string) ([]translate.Word, []interface{}) {
	// initialize the slice of structs that will contain the scraped data
	var words []translate.Word
	var newUrls []interface{}
	// instantiate a new collector object
	c := colly.NewCollector(
		colly.AllowedDomains(domain),
	)
	c.AllowURLRevisit = false

	c.OnHTML(`div.contentListBody.contentLoadMoreSection`, func(el *colly.HTMLElement) {
		el.ForEach("a[href]", func(_ int, e *colly.HTMLElement) {
			url := strings.Replace(e.Attr("href"), "sufinama.org//", "sufinama.org/kalaam/", 1)
			if !slices.Contains(visitedUrls, url) {
				newUrls = append(newUrls, translate.NewUrl{Site: domain, Url: url})
				e.Request.Visit(url)
			}

		})
	})
	// words = scrapeWords(c, words)
	c.OnHTML("div.pMC[data-roman=off]", func(e *colly.HTMLElement) {
		e.ForEach("span[data-m]", func(_ int, el *colly.HTMLElement) {

			// initialize a new Product instance
			word := translate.Word{}

			// scrape the target data
			word.Text = el.Text
			word.Data = el.Attr("data-m")
			word.Occurrence = el.Request.URL.String()

			// add the product instance with scraped data to the list of products
			words = append(words, word)

		})
	})
	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
		// fmt.Println(words)
	})
	for _, url := range urlList {

		c.Visit(url)
	}
	return words, newUrls
}

func scrapeWords(c *colly.Collector, words []translate.Word) []translate.Word {

	c.OnHTML("div.pMC[data-roman=off]", func(e *colly.HTMLElement) {
		e.ForEach("span[data-m]", func(_ int, el *colly.HTMLElement) {

			// initialize a new Product instance
			word := translate.Word{}

			// scrape the target data
			word.Text = el.Text
			word.Data = el.Attr("data-m")

			// add the product instance with scraped data to the list of products
			words = append(words, word)

		})
	})

	return words
}
