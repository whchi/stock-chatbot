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
	"github.com/whchi/stock-chatbot/pkg/gsheet"
	linebotInstance "github.com/whchi/stock-chatbot/pkg/linebot"
)

func LineEventHandler(c *gin.Context) {
	if c.Request.Method != "POST" {
		return
	}
	var Bot = linebotInstance.Setup()
	events, _ := Bot.ParseRequest(c.Request)
	for _, event := range events {
		log.Print(event.Type)
		switch event.Type {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				quota, err := Bot.GetMessageQuota().Do()
				if err != nil {
					log.Println("Quota err:", err)
				}
				log.Println("quota: " + strconv.FormatInt(quota.Value, 10))
				if strings.HasPrefix(message.Text, "ps") {
					stocks := gsheet.FetchData()
					text := message.Text
					replyMsg := template(stocks, text[2:])
					if replyMsg == "" {
						replyMsg = "查無結果"
					}
					if _, err = Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMsg)).Do(); err != nil {
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
	dataLen := len(data)
	if msg == "list" {
		for i := 0; i < dataLen; i++ {
			bgn := data[i]["begin"][0:10]
			end := data[i]["end"][0:10]
			ret += fmt.Sprintf("\n代號: %s, 名稱: %s, 處置期間: %s~%s",
				data[i]["code"], data[i]["name"], bgn, end)
		}
	} else {
		validCode := regexp.MustCompile(`^\d{4,}$`)
		searchKey := "name"
		if validCode.MatchString(msg) {
			searchKey = "code"
		}

		for i := 0; i < dataLen; i++ {
			if data[i][searchKey] == msg {
				bgn := data[i]["begin"][0:10]
				end := data[i]["end"][0:10]
				ret += fmt.Sprintf("代號: %s, 名稱: %s, 處置期間: %s~%s",
					data[i]["code"], data[i]["name"], bgn, end)
				break
			}
		}
	}

	return ret
}
