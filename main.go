package main

import (
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"os"
	"strings"
)

func check(error error) {
	if error != nil {
		fmt.Println(error)
	}
}

func getHtml(url string) *http.Response {
	response, error := http.Get(url)
	check(error)
	if response.StatusCode > 400 {
		fmt.Println("Status code:", response.StatusCode)
	}

	return response
}

func writeCsv(scrapedData []string) {
	filename := "data.csv"
	file, error := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	check(error)
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	error = writer.Write(scrapedData)
	check(error)
}

func scrapePageData(doc *goquery.Document) {
	doc.Find("ul.srp-results>li.s-item").Each(func(index int, item *goquery.Selection) {
		a := item.Find("a.s-item__link")
		title := strings.TrimSpace(a.Text())
		url, _ := a.Attr("href")
		priceSpan := strings.TrimSpace(item.Find("span.s-item__price").Text())
		price := strings.Trim(priceSpan, " руб.")
		scrapedData := []string{title, price, url}
		writeCsv(scrapedData)
	})
}

func main() {
	url := "https://www.ebay.com/sch/i.html?_from=R40&_nkw=board+game+bundle&_sacat=233&LH_TitleDesc=0&_ipg=240"

	var previousUrl string

	for {
		response := getHtml(url)
		defer response.Body.Close()
		doc, error := goquery.NewDocumentFromReader(response.Body)
		check(error)
		scrapePageData(doc)
		href, _ := doc.Find("nav.pagination>a.pagination__next").Attr("href")
		if href == previousUrl {
			break
		} else {
			url = href
			previousUrl = href
		}
	}

}
