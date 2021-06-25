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
				if strings.HasPrefix(message.Text, "/") {
					if !cache.IsExpired() {
						stocks = cache.GetStocks()
					} else {
						stocks = gsheet.FetchData()
						cache.SyncWithRaw(stocks)
					}
					text := message.Text
					replyMsg := template(stocks, text[1:])
					if replyMsg == "" {
						replyMsg = "查無結果"
					}
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

func template(data []map[string]string, msg string) (result string) {
	var ret string

	msg = strings.Trim(msg, " \n")
	dataLen := len(data)

	switch msg {
	case "help":
		ret = "* /list：列出所有處置\n* /股票代碼 or /股票名稱：列出單支股票"
	case "list":
		for i := 0; i < dataLen; i++ {
			bgn := data[i]["begin"][0:10]
			end := data[i]["end"][0:10]
			ret += fmt.Sprintf("代號: %s, 名稱: %s, 處置期間: %s~%s\n",
				data[i]["code"], data[i]["name"], bgn, end)
		}
		ret = ret[:len(ret)-1]
	default:
		validCode := regexp.MustCompile(`^\d{4,}$`)
		searchKey := "name"
		if validCode.MatchString(msg) {
			searchKey = "code"
		}

		for i := 0; i < dataLen; i++ {
			if data[i][searchKey] == msg {
				bgn := data[i]["begin"][0:10]
				end := data[i]["end"][0:10]
				ret += fmt.Sprintf("代號: %s, 名稱: %s, 處置期間: %s~%s\n",
					data[i]["code"], data[i]["name"], bgn, end)
			}
		}
	}

	return ret
}
