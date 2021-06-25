package cache

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func Exists() (exists bool) {
	_, err := os.Stat("data/stocks.json")
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Fatalln(err)
	}
	return true
}

func Sync(data []byte) {
	err := ioutil.WriteFile("data/stocks.json", data, 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func GetStocks() (data []map[string]string) {
	jsonFile, err := os.Open("data/stocks.json")
	if err != nil {
		log.Fatalln(err)
	}

	defer jsonFile.Close()

	var jsonData []map[string]string
	body, _ := ioutil.ReadAll(jsonFile)
	if err := json.Unmarshal(body, &jsonData); err != nil {
		panic(err)
	}

	return jsonData
}
