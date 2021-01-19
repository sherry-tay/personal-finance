FROM golang:1.15

# For local development
# ENV TELEGRAM_BOT_TOKEN="TELEGRAM_BOT_TOKEN"
# ENV TELEGRAM_WEBHOOK_URL="TELEGRAM_WEBHOOK_URL"
# ENV AUTHORIZED_USER="AUTHORIZED_USER"
# ENV PORT="8080"

WORKDIR /app
COPY data ./data
COPY internal ./internal
COPY main.go go.* ./

# For local development
# COPY personal-finance-admin.json ./

RUN go build

EXPOSE 8080
CMD ["/app/personal-finance"]