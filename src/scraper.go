package main

import (
	"fmt"
	"os"
	"sync"

	"sort"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

var wg sync.WaitGroup

func autotrader(collector colly.Collector, make string, model string, year string, c chan []int) []int {
	wg.Add(1)
	defer wg.Done()
	res := []int{}
	myurl := "https://www.autotrader.com/cars-for-sale/all-cars/" + year + "/" + make + "/" + model
	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	collector.OnResponse(func(r *colly.Response) {
		fmt.Println("Got a response from", r.Request.URL)
	})
	collector.OnError(func(r *colly.Response, e error) {
		fmt.Println("An error occurred:", e)
	})
	collector.OnHTML("div", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "first-price") {
			price := e.Text
			toint := string(price[0])
			if strings.Contains(price, ",") {
				for i := 1; i < len(price); i++ {
					if string(price[i]) != "," {
						toint += string(price[i])
					}
				}
			}
			intprice, err := strconv.Atoi(toint)
			if err != nil {
				fmt.Println("fuck")
				fmt.Println(err)
			}
			res = append(res, intprice)
		}
	})
	collector.Visit(myurl)
	sort.Ints(res)
	c <- res
	return res
}

func ebay(collector colly.Collector, make string, model string, year string, c chan []int) []int {
	wg.Add(1)
	defer wg.Done()
	res := []int{}
	myurl := "https://www.ebay.com/sch/i.html?_from=R40&_trksid=p4432023.m570.l1313&_nkw=" + year + "+" + make + "+" + model + "&_sacat=0"
	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	collector.OnResponse(func(r *colly.Response) {
		fmt.Println("Got a response from", r.Request.URL)
	})
	collector.OnError(func(r *colly.Response, e error) {
		fmt.Println("An error occurred:", e)
	})
	collector.OnHTML("span", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "s-item__price") {
			price := e.Text
			// fmt.Println("Price:", price)
			toint := string(price[1])
			if strings.Contains(price, ",") {
				for i := 2; i < len(price); i++ {
					if string(price[i]) == "." {
						break
					}
					if string(price[i]) != "," {
						toint += string(price[i])
					}
				}
			}
			intprice, err := strconv.Atoi(toint)
			if err != nil {
				fmt.Println("fuck")
				fmt.Println(err)
			}
			if intprice > 1000 {
				res = append(res, intprice)
			}
		}
	})
	collector.Visit(myurl)
	sort.Ints(res)
	c <- res
	return res
}

func main() {
	c1 := make(chan []int)
	c2 := make(chan []int)
	args := os.Args
	make := args[1]
	model := args[2]
	var year string
	if len(args) == 4 {
		year = args[3]
	} else {
		year = ""
	}
	collector := colly.NewCollector()
	go autotrader(*collector, make, model, year, c1)
	go ebay(*collector, make, model, year, c2)
	wg.Wait()
	all := append(<-c1, <-c2...)
	sort.Ints(all)
	if len(args) == 4 {
		fmt.Printf("The cheapest %s %s %s costs %d dollars.\n", year, make, model, all[0])
	} else {
		fmt.Printf("The cheapest %s %s costs %d dollars.\n", make, model, all[0])
	}
}
