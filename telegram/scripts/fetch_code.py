#!/usr/bin/env python3

import json
import requests

output_file = 'data/category.json'

def request_all_codes():
    url = 'https://api.sgx.com/securities/v1.1'
    params = {'params' : 'nc,type'}
    r = requests.get(url, params=params)
    return r.json()

with open(output_file, 'w') as outfile:
    dict = {}
    for x in request_all_codes()["data"]["prices"]:
        dict[x["nc"]] = x["type"]
    json.dump(dict, outfile)