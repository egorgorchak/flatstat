package main

import (
	"bytes"
	"flatstat/data"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func main() {
	var url string = "https://www.bazaraki.com/adv/5570371_1-bedroom-apartment-to-rent/"
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	doc, _ := html.Parse(resp.Body)

	getCharacteristics(doc)

	// var id int = 1l
	// var coord string = getCoord(doc)
	// var bedrooms int = getBedrooms(doc)
	// var area int

	// price, _ := strconv.Atoi(getDivContent("announcement-price__cost", doc))

	// var publishDate time.Time
	// var isValid bool

	// var flatId int = 1

	// ad := createAd(flatId, coord, bedrooms, area, id, url, price, publishDate, isValid)

	// fmt.Println(ad)
}

func getDivContent(cl string, doc *html.Node) string {
	for n := range doc.Descendants() {
		if n.DataAtom == atom.Div {
			for _, a := range n.Attr {
				if a.Key == atom.Class.String() && a.Val == cl {
					return strings.TrimSpace(n.LastChild.Data)
				}
			}
		}
	}
	return "0"
}

func getBedrooms(doc *html.Node) int {
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.Div {
			for _, a := range n.Attr {
				if a.Key == atom.Class.String() && a.Val == "announcement-characteristics clearfix" {
					fmt.Println("Bed", getHtmlString(n))
					fmt.Println(n.Type)
					fmt.Println(n.FirstChild.Type)
					return -1
				}
			}
		}
	}
	return -1
}

func getCharacteristics(doc *html.Node) int {
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.Div {
			for _, a := range n.Attr {
				if a.Key == atom.Class.String() && a.Val == "announcement-characteristics clearfix" {
					for d := range n.Descendants() {
						if d.Type == html.ElementNode {
							fmt.Println("Desc", getHtmlString(d))
						}
					}

					return -1
				}
			}
		}
	}
	return -1
}

func getCoord(doc *html.Node) string {
	var lng string
	var lat string
	for n := range doc.Descendants() {
		if n.DataAtom == atom.A {
			for _, a := range n.Attr {
				if a.Key == "data-default-lat" {
					lat = strings.TrimSpace(a.Val)
				}

				if a.Key == "data-default-lng" {
					lng = strings.TrimSpace(a.Val)
				}
			}
		}
	}
	return lat + "," + lng
}

func createAd(flatId int, coord string, bedrooms int, area int, adId int, link string, price int, publishDate time.Time, isValid bool) data.Advertisment {
	return data.Advertisment{
		Id:          adId,
		Link:        link,
		Price:       price,
		PublishDate: publishDate,
		Flat: data.Flat{
			Id:       flatId,
			Coord:    coord,
			Bedrooms: bedrooms,
			Area:     area,
		},
		IsValid: isValid,
	}
}

func getHtmlString(n *html.Node) string {
	var buff bytes.Buffer
	html.Render(&buff, n)

	return buff.String()
}
