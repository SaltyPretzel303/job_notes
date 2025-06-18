import requests
import json
import os

auth_token = os.environ['VH_API_TOKEN']
headers = {'Authorization': 'Token %s' % auth_token}

resp = requests.get('https://app.valohai.com/api/v0/projects/', headers=headers)
resp.raise_for_status()

print(json.dumps(resp.json(), indent=4))



