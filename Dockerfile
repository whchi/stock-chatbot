FROM whcdc/golang1.16-selenium-chrome

WORKDIR /app

COPY . /app

RUN go build -o bin/crawler cmd/crawler/main.go
RUN go build -o bin/webapp cmd/line/main.go

EXPOSE 8080/tcp

CMD ["/app/bin/webapp"]
