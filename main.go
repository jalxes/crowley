package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	u := "https://en.wikipedia.org/wiki/Rider-Waite_tarot_deck"
	fmt.Println(fmt.Sprintf("url: %s", u))

	doc, err := goquery.NewDocument(u)
	if err != nil {
		panic(err)
	}
	type entry struct {
		id    int
		title string
		url   string
		link  string
		desc  string
		cat   interface{}
	}

	var entries []entry
	doc.Find("div.mw-parser-output center ul li").Each(func(i int, node *goquery.Selection) {
		id := i + 1
		title, _ := node.Find("p a").Attr("title")
		title = strings.Replace(title, "(Tarot card)", "", 1)
		title = strings.Replace(title, "(tarot card)", "", 1)

		url, _ := node.Find("a.image").Attr("href")
		link, _ := node.Find("a.image img").Attr("src")
		desc := node.Find("p a").Text()
		// cat, _ := node.ParentFiltered("center").Html()
		link = strings.Replace(link, "thumb/", "", 1)
		link = link[:strings.LastIndex(link, "/")]
		link = "https:" + link
		entries = append(entries, entry{
			id:    id,
			title: title,
			url:   url,
			link:  link,
			desc:  desc,
			// cat:   cat,
		})
	})
	for _, item := range entries {
		fmt.Println(fmt.Sprintf("%d", item.id))
		fmt.Println(fmt.Sprintf("title: %s", item.title))
		fmt.Println(fmt.Sprintf("url: %s", item.url))
		fmt.Println(fmt.Sprintf("link: %s", item.link))
		fmt.Println(fmt.Sprintf("desc: %s", item.desc))
		// fmt.Println(fmt.Sprintf("cat: %s", item.cat))
		fmt.Println("=====================")
		res, err := http.Get(item.link)
		if err != nil {
			log.Fatal(fmt.Sprintf("cant get %s", item.title))
		}
		defer res.Body.Close()
		os.MkdirAll("images", 0777)
		file, err := os.Create(fmt.Sprintf("images/%d-%s.jpg", item.id, item.title))
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		io.Copy(file, res.Body)
	}
}
