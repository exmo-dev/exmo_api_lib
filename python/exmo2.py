import httplib
import urllib
import json
import hashlib
import hmac
import time
 
api_key = "your_key"
api_secret = "your_secret"

nonce = int(round(time.time()*1000))
 
params = {"nonce": nonce}
params = urllib.urlencode(params)
 
H = hmac.new(api_secret, digestmod=hashlib.sha512)
H.update(params)
sign = H.hexdigest()
 
headers = {"Content-type": "application/x-www-form-urlencoded",
           "Key":api_key,
           "Sign":sign}
conn = httplib.HTTPSConnection("api.exmo.com")
conn.request("POST", "/v1/user_info", params, headers)
response = conn.getresponse()
 
print response.status, response.reason
print json.load(response)
 
conn.close()