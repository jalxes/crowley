package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

func main() {
	u := "https://en.wikipedia.org/wiki/Rider-Waite_tarot_deck"
	fmt.Println(fmt.Sprintf("url: %s", u))

	doc, err := htmlquery.LoadURL(u)
	if err != nil {
		panic(err)
	}
	type entry struct {
		id    int
		title string
		url   string
		link  string
		desc  string
	}

	var entries []entry
	htmlquery.FindEach(doc, "//div[@class='mw-parser-output']/center/ul/li", func(i int, node *html.Node) {
		item := entry{}
		item.id = i
		// h2 := htmlquery.FindOne(node, "//h2")
		// item.title = htmlquery.InnerText(h2)
		// item.url = htmlquery.SelectAttr(htmlquery.FindOne(h2, "a"), "href")
		// if n := htmlquery.FindOne(node, "//div[@class='b_caption']/p"); n != nil {
		// 	item.desc = htmlquery.InnerText(n)
		// }
		item.url = htmlquery.SelectAttr(htmlquery.FindOne(node, "//a[@class='image']"), "href")
		item.link = htmlquery.SelectAttr(htmlquery.FindOne(node, "//a[@class='image']/img"), "src")
		item.title = htmlquery.SelectAttr(htmlquery.FindOne(node, "//p/a"), "title")
		item.desc = htmlquery.InnerText(htmlquery.FindOne(node, "//p/a"))
		item.link = strings.Replace(item.link, "thumb/", "", 1)
		item.link = item.link[:strings.LastIndex(item.link, "/")]
		item.link = "https:" + item.link
		entries = append(entries, item)

	})
	for _, item := range entries {
		fmt.Println(fmt.Sprintf("%d", item.id))
		fmt.Println(fmt.Sprintf("%d title: %s", item.id, item.title))
		fmt.Println(fmt.Sprintf("url: %s", item.url))
		fmt.Println(fmt.Sprintf("link: %s", item.link))
		fmt.Println(fmt.Sprintf("desc: %s", item.desc))
		fmt.Println("=====================")
		res, err := http.Get(item.link)
		if err != nil {
			log.Fatal(fmt.Sprintf("cant get %s", item.title))
		}
		defer res.Body.Close()

		file, err := os.Create("tmp/" + item.title + ".jpg")
		if err != nil {
			log.Fatal(fmt.Sprintf("cant create %s", item.title))
		}
		defer file.Close()

		io.Copy(file, res.Body)
	}
}
