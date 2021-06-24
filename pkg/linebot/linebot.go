package linebot

import (
	"log"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/whchi/stock-chatbot/pkg/setting"
)

var Bot *linebot.Client

func Setup() (instance *linebot.Client) {
	Bot, err := linebot.New(
		setting.LineSetting.SECRET,
		setting.LineSetting.ACCESS_TOKEN)

	if err != nil {
		log.Fatal(err)
	}

	return Bot
}
