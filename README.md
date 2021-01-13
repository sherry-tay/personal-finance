# personal-finance

1. [Using Windows Powershell](#using-windows-powershell)
1. [Using MacOS](#using-macOS)
1. [Using Docker](#using-docker)
1. [Create webhook locally using ngrok](#create-webhook-locally-using-ngrok)

## Using Windows Powershell

### To run:

1. Add `TELEGRAM_BOT_TOKEN` and `TELEGRAM_WEBHOOK_URL` as environment variables

1. Update `authorizedUser` variable in `telegram/telegram.go` to Telegram username

1. Run the following command:

```
cd cmd
go build; .\cmd.exe
```

### To run a specific test e.g.:

```
go test .\internal\firestore\
```

### To update SGX securities category info:

```
python .\scripts\fetch_code.py
```

### To update stock holdings in FireStore:

1. Add a JSON file in `data` folder following the format of `stocks.json` in the `templates` directory

1. Run the following:

```
pip3.9 install --upgrade firebase-admin
python3 .\scripts\update_firestore.py
```

Note: `pip3.9 install` only needs to be run if Firebase SDK has not yet been installed

## Using MacOS

### To run:
1. Add `TELEGRAM_BOT_TOKEN` and `TELEGRAM_WEBHOOK_URL` as environment variables

1. Update `authorizedUser` variable in `telegram/telegram.go` to Telegram username

1. Run the following command:

```
cd cmd
go build && ./cmd
```

Or

```
cd cmd
go run main.go
```

### To run a specific test e.g.:

```
go test ./internal/firestore/
```

### To update SGX securities category info:

```
chmod +x scripts/fetch_code.py
scripts/fetch_code.py
```

Note: `chmod +x scripts/fetch_code.py` only needs to be run if permission is denied

### To update stock holdings in FireStore:

1. Add a JSON file in `data` folder following the format of `stocks.json` in the `templates` directory

1. Run the following:

```
chmod +x scripts/update_firestore.py
pip3 install --upgrade firebase-admin
scripts/update_firestore.py
```

Note: `chmod +x scripts/update_firestore.py` and `pip3 install` only needs to be run if permission is denied and Firebase SDK has not yet been installed respectively

## Using Docker

### To run:
1. Update `TELEGRAM_BOT_TOKEN` and `TELEGRAM_WEBHOOK_URL` in `Dockerfile`

1. Update `authorizedUser` variable in `telegram/telegram.go` to Telegram username

1. Run the following command:

```
docker build -t personal-finance .
docker run -it --rm -p 80:80 --name personal-finance personal-finance
```

## Create webhook locally using ngrok

1. Run the below to secure a public URL for port 80 web server

```
ngrok http 80
```

or for Windows Docker with IP `192.168.99.100`:

```
ngrok http 192.168.99.100:80
```

Note: For a web interface for debugging incoming connections, use `http://127.0.0.1:4040`

1. Set https public url as an environment variable

```
export TELEGRAM_WEBHOOK_URL=<url>
```

1. Run Go application as per normal

Note: Use the same terminal session to run the application as the one used to set the environment variable