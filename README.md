# personal-finance

## Using Windows Powershell

To run:
1. Update variables in `telegram/private.go`

1. Run the following command:

```
cd cmd
go build; .\cmd.exe
```

To run a specific test e.g.:

```
go test .\internal\firebase\
```

To update SGX securities category info:

```
python .\scripts\fetch_code.py
```

To update stock holdings in FireStore:

1. Add a JSON file in `data` folder following the format of `stocks.json` in the `templates` directory

1. Run the following:

```
pip3.9 install --upgrade firebase-admin
python3 .\scripts\update_firestore.py
```

Note: `pip3.9 install` only needs to be run if Firebase SDK has not yet been installed

## Using MacOS

To run:
1. Update variables in `telegram/private.go`

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

To run a specific test e.g.:

```
go test ./internal/firebase/
```

To update SGX securities category info:

```
chmod +x scripts/fetch_code.py
scripts/fetch_code.py
```

Note: `chmod +x scripts/fetch_code.py` only needs to be run if permission is denied

To update stock holdings in FireStore:

1. Add a JSON file in `data` folder following the format of `stocks.json` in the `templates` directory

1. Run the following:

```
chmod +x scripts/update_firestore.py
pip3 install --upgrade firebase-admin
scripts/update_firestore.py
```

Note: `chmod +x scripts/update_firestore.py` and `pip3 install` only needs to be run if permission is denied and Firebase SDK has not yet been installed respectively