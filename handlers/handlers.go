package handlers

import (
	"encoding/json"
	"flatstat/data"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Info struct {
}

func (*Info) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if b, err := io.ReadAll(req.Body); err == nil {
		link := string(b)
		getInfo(link, rw)
	}
	defer req.Body.Close()
}

func getInfo(link string, rw http.ResponseWriter) error {
	ad := parsePage(link)
	e := json.NewEncoder(rw)
	return e.Encode(ad)
}

func parsePage(link string) data.Advertisment {
	doc := getDocument(link)

	price, _ := strconv.Atoi(getElementContentByClass(atom.Div, "announcement-price__cost", doc))

	datePosted := strings.Replace(getElementContentByClass(atom.Span, "date-meta", doc), "Posted: ", "", -1)

	publishDate, _ := time.Parse("02.01.2006 15:04", datePosted)

	lat := getAttributeValue(doc, atom.A, "data-default-lat")
	lng := getAttributeValue(doc, atom.A, "data-default-lng")

	bedrooms, _ := strconv.Atoi(getCharacteristic(doc, "Bedrooms:"))
	area := strings.Replace(getCharacteristic(doc, "Property area:"), " m²", "", -1)
	areaConv, _ := strconv.Atoi(area)

	ad := createAd(link, price, publishDate, lat+","+lng, bedrooms, areaConv)

	return ad
}

func getCharacteristic(doc *html.Node, characteristicName string) string {
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.Span {
			if n.FirstChild != nil && n.FirstChild.Data == characteristicName {
				return n.NextSibling.NextSibling.FirstChild.Data
			}
		}
	}
	return ""
}

func getDocument(link string) *html.Node {
	resp, err1 := http.Get(link)

	if err1 != nil {
		fmt.Println("Error getting html: ", err1)
		os.Exit(1)
	}

	defer resp.Body.Close()

	doc, err2 := html.Parse(resp.Body)
	if err2 != nil {
		fmt.Println("Error parsing html: ", err2)
		os.Exit(1)
	}

	return doc
}

func getElementContentByClass(el atom.Atom, cl string, doc *html.Node) string {
	for n := range doc.Descendants() {
		if n.DataAtom == el {
			for _, a := range n.Attr {
				if a.Key == atom.Class.String() && a.Val == cl {
					return strings.TrimSpace(n.LastChild.Data)
				}
			}
		}
	}
	return ""
}

func getAttributeValue(doc *html.Node, at atom.Atom, attrName string) string {
	for n := range doc.Descendants() {
		if n.DataAtom == at {
			for _, a := range n.Attr {
				if a.Key == attrName {
					return strings.TrimSpace(a.Val)
				}
			}
		}
	}
	return ""
}

func createAd(link string, price int, publishDate time.Time, coord string, bedrooms int, area int) data.Advertisment {
	return data.Advertisment{
		Link:        link,
		Price:       price,
		PublishDate: publishDate,
		Flat: data.Flat{
			Coord:    coord,
			Bedrooms: bedrooms,
			Area:     area,
		},
		IsValid: true,
	}
}
