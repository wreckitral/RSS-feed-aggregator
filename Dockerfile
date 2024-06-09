FROM golang:1.22.4-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . ./

RUN go build -o rss-feed .

RUN ls -la

RUN chmod +x ./rss-feed

EXPOSE 7777

CMD ["./rss-feed"]
