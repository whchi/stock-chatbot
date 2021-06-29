package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/whchi/stock-chatbot/pkg/cache"
	"github.com/whchi/stock-chatbot/pkg/setting"
)

var listedData map[string]interface{}

type OtcPrev struct {
	code  string
	name  string
	count string
}

func init() {
	setting.Setup()
}
func main() {
	processPunish()
	processNotice()
	log.Print("all sheet data updated")
}

func processPunish() {
	fileName := "punishing_stocks"
	// process listed
	url := setting.OtherSetting.GSHEET_API_URL + "?tab=" + fileName
	rows, total := getListed("punish")
	var listed []map[string]string
	if total > 0 {
		listed := make([]map[string]string, total)
		for i := range listed {
			listed[i] = make(map[string]string, 7)
		}
		for i := 0; i < total; i++ {
			listed[i]["announce_date"] = convertToEra(rows[i][0], "/") + "T00:00:00+08:00"
			listed[i]["code"] = "'" + rows[i][1]
			listed[i]["name"] = rows[i][2]
			listed[i]["count"] = "'" + rows[i][3]
			interval := strings.Split(rows[i][5], "～")
			listed[i]["begin"] = convertToEra(interval[0], "/") + "T00:00:00+08:00"
			listed[i]["end"] = convertToEra(interval[1], "/") + "T23:59:59+08:00"
		}
	}

	// process otc
	var otc []map[string]string
	doc, recordLen := getOTC(fileName)
	if recordLen > 0 {
		otc = make([]map[string]string, recordLen)
		for i := 0; i < recordLen; i++ {
			otc[i] = make(map[string]string, 7)
		}

		prev := OtcPrev{}
		var code string
		var name string
		var count string
		doc.Find("table#result_table tr:not(:first-child)").Each(func(i int, s *goquery.Selection) {
			if s.Find("td").Length() == 8 {
				code = "'" + s.Find("td:nth-child(3)").Text()
				prev.code = code
				name = s.Find("td:nth-child(4)").Text()
				prev.name = name
				count = "'" + s.Find("td:nth-child(5)").Text()
				prev.count = count
			}
			announceDateAt := 2

			intervalAt := 6

			if s.Find("td").Length() != 8 {
				announceDateAt = 1
				intervalAt = 2
				code = prev.code
				name = prev.name
				count = prev.count
			}
			announce_date := fetchAnnounceDate(s, announceDateAt)
			interval := strings.Split(s.Find("td:nth-child("+strconv.Itoa(intervalAt)+")").Text(), "~")
			begin := convertToEra(interval[0], "/") + "T00:00:00+08:00"
			end := convertToEra(interval[1], "/") + "T23:59:59+08:00"
			otc[i]["announce_date"] = announce_date
			otc[i]["code"] = code
			otc[i]["name"] = name
			otc[i]["count"] = count
			otc[i]["begin"] = begin
			otc[i]["end"] = end
		})
	}

	toDb := append(listed, otc...)
	if len(toDb) > 0 {
		save(url, toDb, fileName+".json")
		log.Print(fileName + " sheet data updated")
	}
}

func processNotice() {
	fileName := "notice_stocks"
	// process listed
	url := setting.OtherSetting.GSHEET_API_URL + "?tab=" + fileName
	rows, total := getListed("notice")
	var listed []map[string]string
	if total > 0 {
		listed = make([]map[string]string, total)
		for i := range listed {
			listed[i] = make(map[string]string, 6)
		}
		for i := 0; i < total; i++ {
			listed[i]["announce_date"] = convertToEra(rows[i][4], ".") + "T00:00:00+08:00"
			listed[i]["code"] = "'" + rows[i][0]
			listed[i]["name"] = rows[i][1]
			listed[i]["count"] = "'" + rows[i][2]
			listed[i]["desc"] = rows[i][3]
		}
	}
	// process otc
	var otc []map[string]string
	doc, recordLen := getOTC(fileName)
	if recordLen > 0 {
		otc = make([]map[string]string, recordLen)
		for i := 0; i < recordLen; i++ {
			otc[i] = make(map[string]string, 6)
		}

		prev := OtcPrev{}
		doc.Find("table#result_table tr:not(:first-child)").Each(func(i int, s *goquery.Selection) {
			descAt := 5
			code := "'" + s.Find("td:nth-child(2)").Text()
			prev.code = code
			name := s.Find("td:nth-child(3)").Text()
			prev.name = name
			count := "'" + s.Find("td:nth-child(4)").Text()
			prev.count = count
			announceDateAt := 6

			if s.Find("td").Length() != 8 {
				announceDateAt = 2
				descAt = 1
				code = prev.code
				name = prev.name
				count = prev.count
			}
			announce_date := fetchAnnounceDate(s, announceDateAt)
			desc := s.Find("td:nth-child(" + strconv.Itoa(descAt) + ")").Text()
			otc[i]["announce_date"] = announce_date
			otc[i]["code"] = code
			otc[i]["name"] = name
			otc[i]["count"] = count
			otc[i]["desc"] = desc
		})
	}

	toDb := append(listed, otc...)
	if len(toDb) > 0 {
		save(url, toDb, fileName+".json")
		log.Print(fileName + " sheet data updated")
	}
}

func fetchAnnounceDate(s *goquery.Selection, at int) (rst string) {
	return convertToEra(s.Find("td:nth-child("+strconv.Itoa(at)+")").Text(), "/") + "T00:00:00+08:00"
}

func getOTC(fileName string) (raw *goquery.Document, total int) {
	urlSuffix := "/disposal_information/disposal_information.php?l=zh-tw"
	if fileName == "notice_stocks" {
		urlSuffix = "/attention_information/trading_attention_information.php?l=zh-tw"
	}
	url := "https://www.tpex.org.tw/web/bulletin" + urlSuffix
	rawHtml := SeleniumFetch(url)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawHtml))
	if err != nil {
		log.Fatalln(err)
	}
	recordLen := doc.Find("table#result_table tr:not(:first-child)").Length()

	return doc, recordLen
}

func convertToEra(old string, sep string) (new string) {
	parsed := strings.Split(old, sep)
	yr, _ := strconv.Atoi(parsed[0])
	yr += 1911

	return strconv.Itoa(yr) + "-" + parsed[1] + "-" + parsed[2]
}
func getListed(stockType string) (results [][]string, count int) {
	t := time.Now()
	startDate := t.Format("20060102")
	endDate := t.AddDate(0, 0, 1).Format("20060102")
	ts := strconv.Itoa(int(time.Now().Unix()))
	url := fmt.Sprintf("https://www.twse.com.tw/announcement/%s?response=json&startDate=%s&endDate=%s&radioType=false&_=%s", stockType, startDate, endDate, ts)

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
				for i := range r {
					if i > 0 && i < 7 {
						rec := fmt.Sprintf("%v", r[i])
						rows[idx][j] = rec
						j += 1
					}
				}
			}
		}
	}

	return rows, total
}

func SeleniumFetch(url string) (result string) {
	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.NoSandbox,
	}
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
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
func save(url string, data []map[string]string, fileName string) bool {
	if payload, err := json.Marshal(data); err == nil {
		_, err := http.Post(url, "Application/json", strings.NewReader(string(payload)))
		if err != nil {
			log.Fatalln(err)
		}
		cache.Sync(payload, fileName)
		return true
	}
	return false
}
