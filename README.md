[![CircleCI](https://circleci.com/gh/sherry-tay/personal-finance.svg?style=svg&circle-token=24d608af19fa19829456e62b759188a6a6273112)](https://app.circleci.com/pipelines/github/sherry-tay/personal-finance)

# personal-finance

1. [Using Windows Powershell](#using-windows-powershell)
1. [Using MacOS](#using-macOS)
1. [Using Docker](#using-docker)
1. [Create webhook locally using ngrok](#create-webhook-locally-using-ngrok)

## Using Windows Powershell

### To run

1. Add `TELEGRAM_BOT_TOKEN`, `TELEGRAM_WEBHOOK_URL` and `AUTHORIZED_USER` (Telegram username) as environment variables

1. Run the following command:

```
go build; .\personal-finance.exe
```

### To run a specific test

```
go test .\internal\firestore\
```

### To update SGX securities category info

```
python .\scripts\fetch_code.py
```

### To update stock holdings in FireStore

1. Add a JSON file in `data` folder following the format of `stocks.json` in the `templates` directory

1. Run the following:

```
pip3.9 install --upgrade firebase-admin
python3 .\scripts\update_firestore.py
```

Note: `pip3.9 install` only needs to be run if Firebase SDK has not yet been installed

## Using MacOS

### To run

1. Add `TELEGRAM_BOT_TOKEN`, `TELEGRAM_WEBHOOK_URL` and `AUTHORIZED_USER` (Telegram username) as environment variables

1. Run the following command:

```
go build && ./personal-finance
```

Or

```
go run main.go
```

### To run a specific test

```
go test ./internal/firestore/
```

### To run all tests including subdirectories

```
go test ./...
```

### To update SGX securities category info

```
chmod +x scripts/fetch_code.py
scripts/fetch_code.py
```

Note: `chmod +x scripts/fetch_code.py` only needs to be run if permission is denied

### To update stock holdings in FireStore

1. Add a JSON file in `data` folder following the format of `stocks.json` in the `templates` directory

1. Run the following:

```
chmod +x scripts/update_firestore.py
pip3 install --upgrade firebase-admin
scripts/update_firestore.py
```

Note: `chmod +x scripts/update_firestore.py` and `pip3 install` only needs to be run if permission is denied and Firebase SDK has not yet been installed respectively

## Using Docker

### To run

1. Update `TELEGRAM_BOT_TOKEN`, `TELEGRAM_WEBHOOK_URL` and `AUTHORIZED_USER` (Telegram username) in `Dockerfile`

1. Run the following command:

```
docker build -t personal-finance .
docker run -it --rm -p 8080:8080 --name personal-finance personal-finance
```

## Create webhook locally using ngrok

1. Run the below to secure a public URL for port 8080 web server

```
ngrok http 8080
```

or for Windows Docker with IP `192.168.99.100`:

```
ngrok http 192.168.99.100:8080
```

Note: For a web interface for debugging incoming connections, use `http://127.0.0.1:4040`

1. Set https public url as an environment variable

```
export TELEGRAM_WEBHOOK_URL=<url>
```

1. Run Go application as per normal

Note: Use the same terminal session to run the application as the one used to set the environment variable