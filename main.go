
package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
)

func main() {
	allowedChannels := []string{"TOKYO MX", "BS11", "日本テレビ", "フジテレビ"}
	url := "https://akiba-souken.com/anime/autumn/"
	doc, _ := goquery.NewDocument(url)

	doc.Find("div.main div.itemBox").Each(func(i int, s *goquery.Selection){
	//	text := s.Find("div.mTitle h2").Text()

		station := make(map[int]map[string]string)
		s.Find("div.itemData div.schedule table tbody tr td span.station").Each(func(j int, s *goquery.Selection) {
			for _, channel := range allowedChannels {
				broadcast := s.Text()
				Time := s.Next().Text()
				if channel == broadcast {
					station[j] = make(map[string]string)
					station[j][broadcast] = strings.Replace(Time, "～", "", -1)
				}
			}
		})

		var ranking []string
		s.Find("div.itemData div.related ul.link li a").Each(func(_ int, s *goquery.Selection) {
			//log.Println(s.Html())
			ranking = append(ranking, s.Text())
		})

		// アニメタイトル
		//log.Println("TITLE : " + text)
		// アニメランキング
		//log.Println(ranking)
		// 放送局
		for _, s := range station {
			for _, k := range s {
					time := "00:00"
					t := k
					split := strings.Split(t, "年")
					year := split[0]
					split = strings.Split(split[1], "月")
					month := split[0]
					split = strings.Split(split[1], "日")
					day := split[0]
					log.Println(t)
					split = strings.Split(split[1], ")")
					if len(split) > 1 {
						time = split[1]
					}

					log.Println(year + "-" + month + "-" + day + "T" + time+ "+09:00")
			}

		}
		log.Println("===============")
	})

}