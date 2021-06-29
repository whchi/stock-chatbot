package api

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/whchi/stock-chatbot/pkg/cache"
	"github.com/whchi/stock-chatbot/pkg/gsheet"
	linebotInstance "github.com/whchi/stock-chatbot/pkg/linebot"
)

func LineEventHandler(c *gin.Context) {
	if c.Request.Method != "POST" {
		return
	}
	var Bot = linebotInstance.Setup()
	events, _ := Bot.ParseRequest(c.Request)
	var stocks []map[string]string
	for _, event := range events {
		log.Print(event.Type)
		switch event.Type {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if message.Text == "how do you turn this on" {
					cache.Flush()
					if _, err := Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("嗚痾")).Do(); err != nil {
						log.Panic(err)
					}
					return
				}
				if strings.HasPrefix(message.Text, "/") {
					fileName := "punishing_stocks.json"
					if !cache.IsExpired(fileName) {
						stocks = cache.GetStocks(fileName)
					} else {
						stocks = gsheet.FetchData()
						cache.SyncWithRaw(stocks, fileName)
					}
					text := message.Text
					search := text[1:]
					if len(search) == 0 {
						if _, err := Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("打字好ㄇ")).Do(); err != nil {
							log.Panic(err)
						}
						return
					}
					replyMsg := template(stocks, search, fileName)
					if _, err := Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMsg)).Do(); err != nil {
						log.Panic(err)
					}
				} else if strings.HasPrefix(message.Text, "!") {
					fileName := "notice_stocks.json"
					if !cache.IsExpired(fileName) {
						stocks = cache.GetStocks(fileName)
					} else {
						stocks = gsheet.FetchData()
						cache.SyncWithRaw(stocks, fileName)
					}
					text := message.Text
					search := text[1:]
					if len(search) == 0 {
						if _, err := Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("打字好ㄇ")).Do(); err != nil {
							log.Panic(err)
						}
						return
					}
					replyMsg := template(stocks, search, fileName)
					if _, err := Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMsg)).Do(); err != nil {
						log.Panic(err)
					}
				} else if strings.HasPrefix(message.Text, "!") {
					fileName := "notice_stocks.json"
					if !cache.IsExpired(fileName) {
						stocks = cache.GetStocks(fileName)
					} else {
						stocks = gsheet.FetchData()
						cache.SyncWithRaw(stocks, fileName)
					}
					text := message.Text
					search := text[1:]
					if len(search) == 0 {
						if _, err := Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("打字好ㄇ")).Do(); err != nil {
							log.Panic(err)
						}
						return
					}
					replyMsg := template(stocks, search, fileName)
					if _, err := Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMsg)).Do(); err != nil {
						log.Panic(err)
					}
				} else {
					rand.Seed(time.Now().Unix())
					number := strconv.Itoa(rand.Intn(9999))
					if _, err := Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("歐印 "+number+"，財富自由")).Do(); err != nil {
						log.Panic(err)
					}
				}
			default:
				if _, err := Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("致富密碼 2330，一張換一個功能")).Do(); err != nil {
					log.Panic(err)
				}
			}
		default:
			Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("致富密碼 2330，一張換一個功能"))
		}
	}
}

func template(data []map[string]string, msg string, fileName string) (result string) {
	var ret string

	msg = strings.Trim(msg, " \n")
	dataLen := len(data)

	switch msg {
	case "help":
		ret = "使用 '/' 查詢處置股票，'!' 查詢注意股票\n\n* /list：列出所有處置\n* (/ or !){股票代碼} or (/ or !){股票名稱}：列出單支股票\n\n範例：\n!台積電\n/2330"
	case "list":
		if fileName != "punishing_stocks.json" {
			ret = "注意股票字太多了，清單會爆"
		} else {
			for i := 0; i < dataLen; i++ {
				bgn := data[i]["begin"][0:10]
				end := data[i]["end"][0:10]
				ret += fmt.Sprintf("代號: %s, 名稱: %s, 處置期間: %s~%s\n",
					data[i]["code"][1:], data[i]["name"], bgn, end)
			}
			ret = ret[:len(ret)-1]
		}
	default:
		validCode := regexp.MustCompile(`^\d{4,}$`)
		searchKey := "name"
		if validCode.MatchString(msg) {
			searchKey = "code"
			msg = "'" + msg
		}
		for i := 0; i < dataLen; i++ {
			fmt.Println(data[i][searchKey], msg)
			if data[i][searchKey] == msg {
				if fileName != "punishing_stocks.json" {
					ret += fmt.Sprintf("-----\n代號: %s\n名稱: %s\n理由: %s\n,宣布日期: %s\n",
						data[i]["code"][1:], data[i]["name"], data[i]["desc"], data[i]["announce_date"])
				} else {
					bgn := data[i]["begin"][0:10]
					end := data[i]["end"][0:10]
					ret += fmt.Sprintf("代號: %s, 名稱: %s, 處置期間: %s~%s\n",
						data[i]["code"][1:], data[i]["name"], bgn, end)
				}
			}
		}
	}

	if ret == "" {
		ret = "查無結果"
	}
	return ret
}
