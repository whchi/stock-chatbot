package gsheet

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/whchi/stock-chatbot/pkg/setting"
)

var r []map[string]string

func FetchData(tab string) (result []map[string]string) {
	url := setting.OtherSetting.GSHEET_API_URL + "?tab=" + tab
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
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if err := json.Unmarshal(body, &r); err != nil {
		panic(err)
	}

	return r
}
