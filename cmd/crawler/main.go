package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/joho/godotenv"
)

var listedData map[string]interface{}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	url := os.Getenv("GSHEET_API_URL") + "?tab=punishing_stocks"
	listed := GetListed()
	otc := GetOTC()
	toDb := append(listed, otc...)
	SaveToSheet(url, toDb)
}

type OtcPrev struct {
	code         string
	name         string
	punish_count string
}

func GetOTC() (toDb []map[string]string) {
	t := time.Now()
	startDate := t.Format("2006/01/02")
	endDate := t.AddDate(0, 0, 1).Format("2006/01/02")

	payload := url.Values{}
	payload.Set("input_date_Start", string(startDate))
	payload.Set("input_date_End", string(endDate))
	payload.Set("Submit", "查詢")
	payload.Set("choice_type", "all_category")
	payload.Set("group_type", "group_stk")
	payload.Set("disposal_reason", "-1")
	payload.Set("disposal_measure", "-1")
	payload.Set("stk_code", "")
	payload.Set("stk_cotegory", "-1")
	url := "https://www.tpex.org.tw/web/bulletin/disposal_information/disposal_information.php?l=zh-tw"
	rawHtml := SeleniumFetch(url)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawHtml))
	if err != nil {
		log.Fatalln(err)
	}
	recordLen := doc.Find("table#result_table tr:not(:first-child)").Length()
	result := make([]map[string]string, recordLen)

	for i := 0; i < recordLen; i++ {
		result[i] = make(map[string]string, 7)
	}

	prev := OtcPrev{}
	doc.Find("table#result_table tr:not(:first-child)").Each(func(i int, s *goquery.Selection) {
		var code string
		var name string
		var punish_count string
		if s.Find("td").Length() == 8 {
			code = s.Find("td:nth-child(3)").Text()
			prev.code = code
			name = s.Find("td:nth-child(4)").Text()
			prev.name = name
			punish_count = s.Find("td:nth-child(5)").Text()
			prev.punish_count = punish_count
		}
		announceDateAt := 2

		intervalAt := 6

		if s.Find("td").Length() != 8 {
			announceDateAt = 1
			intervalAt = 2
			code = prev.code
			name = prev.name
			punish_count = prev.punish_count
		}
		announce_date := convertToEra(s.Find("td:nth-child("+strconv.Itoa(announceDateAt)+")").Text()) + "T00:00:00+08:00"
		interval := strings.Split(s.Find("td:nth-child("+strconv.Itoa(intervalAt)+")").Text(), "~")
		begin := convertToEra(interval[0]) + "T00:00:00+08:00"
		end := convertToEra(interval[1]) + "T23:59:59+08:00"
		result[i]["announce_date"] = announce_date
		result[i]["code"] = code
		result[i]["name"] = name
		result[i]["punish_count"] = punish_count
		result[i]["begin"] = begin
		result[i]["end"] = end
	})

	return result
}

func convertToEra(old string) (new string) {
	parsed := strings.Split(old, "/")
	yr, _ := strconv.Atoi(parsed[0])
	yr += 1911

	return strconv.Itoa(yr) + "-" + parsed[1] + "-" + parsed[2]
}
func GetListed() (result []map[string]string) {
	t := time.Now()
	startDate := t.Format("20060102")
	endDate := t.AddDate(0, 0, 1).Format("20060102")
	ts := strconv.Itoa(int(time.Now().Unix()))
	url := fmt.Sprintf("https://www.twse.com.tw/announcement/punish?response=json&startDate=%s&endDate=%s&radioType=false&_=%s", startDate, endDate, ts)

	resBody := Fetch(url)
	body, err := ioutil.ReadAll(resBody)
	if err != nil {
		log.Fatalln(err)
	}
	defer resBody.Close()

	if err := json.Unmarshal(body, &listedData); err != nil {
		panic(err)
	}

	total := int(listedData["total"].(float64))
	rows := make([][]string, total)
	for i := range rows {
		rows[i] = make([]string, 6)
	}
	if v, ok := listedData["data"].([]interface{}); ok {
		for idx, row := range v {
			j := 0
			if r, yes := row.([]interface{}); yes {
				for i, _ := range r {
					if i > 0 && i < 7 {
						rec := fmt.Sprintf("%v", r[i])
						rows[idx][j] = rec
						j += 1
					}
				}
			}
		}
	}

	toDb := make([]map[string]string, total)
	for i := range toDb {
		toDb[i] = make(map[string]string, 7)
	}
	for i := 0; i < total; i++ {
		toDb[i]["announce_date"] = rows[i][0] + "T00:00:00+08:00"
		toDb[i]["code"] = rows[i][1]
		toDb[i]["name"] = rows[i][2]
		toDb[i]["punish_count"] = rows[i][3]
		interval := strings.Split(rows[i][5], "～")
		toDb[i]["begin"] = interval[0] + "T00:00:00+08:00"
		toDb[i]["end"] = interval[1] + "T23:59:59+08:00"
	}

	return toDb
}

func SeleniumFetch(url string) (result string) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var htmlContent string

	err := chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate(url),
		chromedp.WaitVisible(`#result_table`, chromedp.ByID),
		chromedp.OuterHTML("html", &htmlContent))
	if err != nil {
		log.Fatalln(err)
	}
	return htmlContent
}

func Fetch(url string) (result io.ReadCloser) {
	client := http.DefaultClient
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X x.y; rv:42.0) Gecko/20100101 Firefox/42.0")
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	return res.Body
}

// [
//     {
//         "code": "2603",
//         "begin": "2021-01-01",
//         "end": "2021-01-02",
//         "name": "長榮",
//         "punish_count":"1",
//         "announce_date": "2020-01-01"
//     }
// ]
func SaveToSheet(url string, data []map[string]string) bool {
	if payload, err := json.Marshal(data); err == nil {
		_, err := http.Post(url, "Application/json", strings.NewReader(string(payload)))
		if err != nil {
			log.Fatalln(err)
		}
		return true
	}
	return false
}
