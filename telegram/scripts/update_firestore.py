#!/usr/bin/env python3

import json
import datetime
import firebase_admin
from firebase_admin import credentials
from firebase_admin import firestore

serviceAccountFilePath = 'internal/firestore/personal-finance-admin.json'
input_file = 'data/stocks.json'
stocksCollectionName = 'stocks'

cred = credentials.Certificate(serviceAccountFilePath)
firebase_admin.initialize_app(cred)

db = firestore.client()

with open(input_file) as json_file:
    data = json.load(json_file)
    for entry in data:
        date = entry["date"]
        entry["date"] = datetime.datetime(int(date[0:4]), int(date[4:6]), int(date[6:8]), tzinfo=datetime.timezone(datetime.timedelta(hours=8)))
        documentId = date + "-" + entry["code"]
        db.collection(stocksCollectionName).document(documentId).set(entry)