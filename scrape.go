package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

type Card struct {
	productName string
	rarity string
	number string
	marketPrice string
	listedMedian string
}

// main() contains code adapted from example found in Colly's docs:
// http://go-colly.org/docs/examples/basic/
func main() {
	// Instantiate default collector
	c := colly.NewCollector()
	foundCardList := []Card{}

	// On every a element which has href attribute call callback
	c.OnHTML(".priceGuideTable", func(e *colly.HTMLElement) {

		e.ForEach(".priceGuideTable tbody tr", func(i int, e *colly.HTMLElement) {
			fmt.Println("Line: ", i)
			var foundCard Card
			foundCard.productName = e.ChildText(".product .cellWrapper .productDetail a")
			foundCard.rarity = e.ChildText(".rarity .cellWrapper")
			foundCard.number = e.ChildText(".number .cellWrapper")
			foundCard.marketPrice = e.ChildText(".marketPrice .cellWrapper")
			foundCard.listedMedian = e.ChildText(".medianPrice .cellWrapper")
			// fmt.Println("Card Info: ", foundCard)
			foundCardList = append(foundCardList, foundCard)
		})

		// Print link
		fmt.Printf("Link Table: %v\n", e.Name)
	})

	c.OnError(func(_ *colly.Response, err error) {
        fmt.Println("Something went wrong:", err)
    })

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnScraped(func(r *colly.Response) {
        fmt.Println("Finished", r.Request.URL)

		fmt.Println("FoundCardList: ", foundCardList)
    })

	// Start scraping on https://hackerspaces.org
	c.Visit("https://shop.tcgplayer.com/price-guide/pokemon")
}
