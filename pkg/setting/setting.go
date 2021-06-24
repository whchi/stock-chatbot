package setting

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Server struct {
	RunMode string
}

type Line struct {
	ACCESS_TOKEN string
	SECRET       string
}

type Other struct {
	GSHEET_API_URL string
}

var OtherSetting = &Other{}
var ServerSetting = &Server{}
var LineSetting = &Line{}

func Setup() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("setting.Setup, fail to get .env: %v", err)
	}
	ServerSetting.RunMode = os.Getenv("RUN_MODE")
	LineSetting.ACCESS_TOKEN = os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	LineSetting.SECRET = os.Getenv("LINE_CHANNEL_SECRET")
	OtherSetting.GSHEET_API_URL = os.Getenv("GSHEET_API_URL")
}
