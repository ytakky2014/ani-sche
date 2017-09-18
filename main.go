
package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
	"time"
	"fmt"
	"strconv"
	"./mygoogle"
)


func main() {
	allowedChannels := []string{"TOKYO MX", "TBS", "テレビ東京", "日本テレビ", "フジテレビ", "BS11", "BS-TBS"}
	url := "https://akiba-souken.com/anime/autumn/"
	doc, _ := goquery.NewDocument(url)

	animeIn := mygoogle.Anime{}
	doc.Find("div.main div.itemBox").Each(func(i int, s *goquery.Selection) {
		// Title
		animeIn.Title = s.Find("div.mTitle h2").Text()

		// Ranking
		var ranking []string
		s.Find("div.itemData div.related ul.link li a").Each(func(_ int, s *goquery.Selection) {
			ranking = append(ranking, s.Text())
		})
		animeIn.Ranking = ranking

		log.Println(animeIn.Title)
		// 放送局
		s.Find("div.itemData div.schedule table tbody tr td span.station").EachWithBreak(func(j int, s *goquery.Selection) bool{
			// 指定放送局のみ視聴する
			for _, channel := range allowedChannels {
				broadcast := s.Text()
				Time := s.Next().Text()
				if channel == broadcast {
					animeIn.Station = broadcast
					time := strings.Replace(Time, "～", "", -1)
					startTime, endTime := convertTime(time)
					animeIn.StartTime = startTime
					animeIn.EndTime = endTime
					log.Println("開始時間" + startTime)
					log.Println("終了時間" + endTime)
					return false
				}
			}
			return true
		})
		out := mygoogle.CreateCalender(animeIn)
		log.Println(out)

	})

}

func convertTime(animeTime string) (string, string) {
	startTime := "00:00"
    addTime := 0
	split := strings.Split(animeTime, "年")
	year, _  := strconv.Atoi(split[0])
	split = strings.Split(split[1], "月")
	month, _ :=  strconv.Atoi(split[0])
	split = strings.Split(split[1], "日")
	day , _ :=  strconv.Atoi(split[0])
	split = strings.Split(split[1], ")")
	h := 0
	m := 0
	if len(split) > 1  && split[1] != "" {
		startTime = split[1]
		// 放送日時が確定している場合は30分で枠を取る
		addTime = 30
		split = strings.Split(startTime, ":")
		h, _ = strconv.Atoi(split[0])
		m, _  = strconv.Atoi(split[1])
		// 24時以降ならば日付を1日増やして24時間減算する
		if h >= 24 {
			day = day + 1
			h = h - 24
		}
	}

	// 各数値をstring化
	yearString := fmt.Sprintf("%04s", strconv.Itoa(year))
	monthString := fmt.Sprintf("%02s", strconv.Itoa(month))
	dayString := fmt.Sprintf("%02s", strconv.Itoa(day))
	hourString := fmt.Sprintf("%02s", strconv.Itoa(h))
	minuteString := fmt.Sprintf("%02s", strconv.Itoa(m))
	startTimeString := yearString + "-" + monthString + "-" + dayString + "T" + hourString + ":" + minuteString + ":00+09:00"
	layout := "2006-01-02T15:04:05-07:00"
	t, _ := time.Parse(layout, startTimeString)
	t = t.Add(time.Duration(addTime) * time.Minute)
	return startTimeString , t.Format(layout)
}