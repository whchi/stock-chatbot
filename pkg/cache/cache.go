package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func Sync(data []byte, fileName string) {
	err := ioutil.WriteFile("data/"+fileName, data, 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func IsExpired(fileName string) (exists bool) {
	jsonFile, err := os.Stat("data/" + fileName)
	if err != nil {
		return true
	}
	execAt := time.Now().Unix()
	lastModifiedAt := jsonFile.ModTime().Unix()

	return execAt-lastModifiedAt > 86400
}

func SyncWithRaw(data []map[string]string, fileName string) {
	if payload, err := json.Marshal(data); err == nil {
		Sync(payload, fileName)
	}
}

func GetStocks(fileName string) (data []map[string]string) {
	jsonFile, err := os.Open("data/" + fileName)
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

func Flush() {
	files, _ := ioutil.ReadDir("data")
	for _, f := range files {
		fmt.Println(f.Name())
		if f.Name() == ".gitkeep" {
			continue
		}
		os.Remove("data/" + f.Name())
	}
}
