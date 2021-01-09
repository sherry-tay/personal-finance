FROM golang:1.15

ENV TELEGRAM_BOT_TOKEN="TELEGRAM_BOT_TOKEN"
ENV TELEGRAM_WEBHOOK_URL="TELEGRAM_WEBHOOK_URL"

WORKDIR /go/src/app
COPY cmd ./cmd
COPY data ./data
COPY internal ./internal
COPY go.* ./

WORKDIR /go/src/app/cmd
RUN go build

EXPOSE 80
CMD ["/go/src/app/cmd/cmd"]