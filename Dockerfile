FROM whcdc/golang1.16-selenium-chrome

RUN apt install -y supervisor

COPY deployments/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

WORKDIR /app

COPY . /app

RUN go build -o bin/crawler cmd/crawler/main.go
RUN go build -o bin/webapp cmd/line/main.go

ENTRYPOINT ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
