# stock-chatbot
golang based side project

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
go run chatbot/cmd/crawler
```
