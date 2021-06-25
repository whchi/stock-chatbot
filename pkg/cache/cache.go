package cache

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func Sync(data []byte) {
	err := ioutil.WriteFile("data/stocks.json", data, 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func IsExpired() (exists bool) {
	jsonFile, err := os.Stat("data/stocks.json")
	if err != nil {
		return true
	}
	execAt := time.Now().Unix()
	lastModifiedAt := jsonFile.ModTime().Unix()

	return execAt-lastModifiedAt > 86400
}

func SyncWithRaw(data []map[string]string) {
	if payload, err := json.Marshal(data); err == nil {
		Sync(payload)
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
