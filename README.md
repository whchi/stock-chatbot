# stock-chatbot
golang based side project
# data source(daily update)
* 上市

[處置](https://www.twse.com.tw/zh/page/announcement/punish.html)

[注意](https://www.twse.com.tw/zh/page/announcement/notice.html)
* 上櫃

[處置](https://www.tpex.org.tw/web/bulletin/disposal_information/disposal_information.php?l=zh-tw)

[注意](https://www.tpex.org.tw/web/bulletin/attention_information/trading_attention_information.php?l=zh-tw)
# develop
```sh
docker-compose up -d
docker exec -ti stock-chatbot_app_1 bash
# in container
go mod tidy
```
# cmd
* fetch data (in container)
```sh
go run cmd/crawler/main.go
```
📈📉📈📉📈📉📈📉📈📉📈📉📈📉📈📉📈📉📈📉📈📉

![](http://shachiku.life/extra/stock-chatbot.png)

📉📈📉📈📉📈📉📈📉📈📉📈📉📈📉📈📉📈📉📈📉📈
