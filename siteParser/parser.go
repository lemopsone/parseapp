package siteParser

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/lemopsone/parseapp/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func getBrandNames() []string {
	return []string{"Bape", "Undercover", "Yohji+Yamamoto",
		"Cav+Empt", "Sacai", "Number+Nine", "Takahiromiyashita+TheSoloist",
		"Mastermind+Japan", "Visvim", "Issey+Miyake", "Neighbourhood"}
}

func Parse() models.ItemList {
	// current stats: approx 40k items get fetched in 2:40 (rough estimate)
	// TO-DO: add database integration
	var wg sync.WaitGroup
	var itemList models.ItemList

	for _, brandName := range getBrandNames() {
		wg.Add(1)
		go retrieveAllPagesForBrand(brandName, &wg, &itemList)
	}
	wg.Wait()
	fmt.Printf("%d items parsed", len(itemList.Items))
	return itemList
}

func retrieveItemCurrentPrice(s *goquery.Selection) (string, string) {
	currentPrice := s.Find("td.currentPrice a").Text()
	if currentPrice == "-" {
		return "-", "-"
	}
	return currentPrice, s.Find("td.currentPrice small.priceRUR").Text()
}

func retrieveItemBlitzPrice(s *goquery.Selection) (string, string) {
	blitzPrice := s.Find("td.bidOrBuy").Text()
	if blitzPrice == "-" {
		return "-", "-"
	}
	splitPrices := strings.Split(blitzPrice, "\n")
	return splitPrices[0], splitPrices[1]
}

func printItemDetails(item models.Item) error {
	b, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func retrievePageItems(url string, itemList *models.ItemList, wg *sync.WaitGroup) {
	defer wg.Done()
	page, _ := loadPage(url)
	page.Find("table.itemList tr").Each(func(i int, s *goquery.Selection) {
		title := s.Find("a.itemName").Text()
		href, _ := s.Find("a.itemName").Attr("href")
		item := models.Item{}
		if title != "" {
			// fmt.Printf("Item %d: %s\n", i, title)
			item.Href = href
			item.Title = title
			item.CurrentPriceYEN, item.CurrentPriceRUR = retrieveItemCurrentPrice(s)
			item.BlitzPriceYEN, item.BlitzPriceRUR = retrieveItemBlitzPrice(s)
			item.TimeLeft = s.Find("td.timeLeft").Text()
			item.TelegramID = ""
			itemList.Items = append(itemList.Items, item)
		}
	})
}

func retrieveAllPagesForBrand(brandName string, wg *sync.WaitGroup, itemList *models.ItemList) {
	defer wg.Done()
	page, _ := loadPage(fmt.Sprintf("https://www.bestjapan.ru/auction/search?page=1&part_id=23000&q=%s",
		brandName))
	totalString := page.
		Find("table.catNavigator tbody tr th").
		Last().Text()
	itemsCount, _ := strconv.Atoi(strings.Split(totalString, " ")[1])
	pagesCount := itemsCount / 20
	if itemsCount%20 != 0 {
		pagesCount += 1
	}
	for i := 1; i <= pagesCount+3; i++ {
		url := fmt.Sprintf("https://www.bestjapan.ru/auction/search?page=%d&part_id=23000&q=%s",
			i, brandName)
		wg.Add(1)
		go retrievePageItems(url, &*itemList, &*wg)
		time.Sleep(500 * time.Millisecond)
	}
}

func loadPage(requestURL string) (*goquery.Document, error) {
	log.Println(requestURL)
	res, err := http.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	return goquery.NewDocumentFromReader(res.Body)
}
