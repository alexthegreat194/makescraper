package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocolly/colly"
)

type Card struct {
	ProductName string `json:"productName"`
	Rarity string `json:"rarity"`
	Number string `json:"number"`
	MarketPrice string `json:"marketPrice"`
	ListedMedian string `json:"listedMedian"`
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
			// fmt.Println("Line: ", i)
			var foundCard Card
			foundCard.ProductName = e.ChildText(".product .cellWrapper .productDetail a")
			foundCard.Rarity = e.ChildText(".rarity .cellWrapper")
			foundCard.Number = e.ChildText(".number .cellWrapper")
			foundCard.MarketPrice = e.ChildText(".marketPrice .cellWrapper")
			foundCard.ListedMedian = e.ChildText(".medianPrice .cellWrapper")
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
		jsonList := []string{}
		for _, v := range foundCardList {
			data, err := json.Marshal(v)
			if (err != nil) {
				panic(err)
			}
			// fmt.Println(v, string(data))
			jsonList = append(jsonList, string(data))
		}

		startString := "{\"url\":\"https://shop.tcgplayer.com/price-guide/pokemon\", \"data\":["
		endString := "]}"

		f, err := os.Create("output.json")
		if (err != nil) {
			panic(err)
		}

		fs := bufio.NewWriter(f)
		fs.WriteString(startString)
		for i, v := range jsonList {
			b, err := fs.WriteString(v)
			if (err != nil) {
				panic(err)
			}
			fmt.Printf("Added %v : %v : %v to File\n", i, v, b)
		}
		fs.WriteString(endString)

		fmt.Printf("Added to File\n")
    })

	// Start scraping on https://hackerspaces.org
	c.Visit("https://shop.tcgplayer.com/price-guide/pokemon")
}
