FROM golang:1.15

ENV TELEGRAM_BOT_TOKEN="TELEGRAM_BOT_TOKEN"
ENV TELEGRAM_WEBHOOK_URL="TELEGRAM_WEBHOOK_URL"
ENV AUTHORIZED_USER="AUTHORIZED_USER"

WORKDIR /go/src/app
COPY data ./data
COPY internal ./internal
COPY main.go go.* ./

RUN go build

EXPOSE 80
CMD ["/go/src/app/personal-finance"]