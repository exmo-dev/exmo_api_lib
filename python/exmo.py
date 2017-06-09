import http.client
import urllib
import json
import hashlib
import hmac
import time

class ExmoAPI:
    def __init__(self, API_KEY, API_SECRET, API_URL = 'api.exmo.com', API_VERSION = 'v1'):
        self.API_URL = API_URL
        self.API_VERSION = API_VERSION
        self.API_KEY = API_KEY
        self.API_SECRET = bytes(API_SECRET, encoding='utf-8')

    def sha512(self, data):
        H = hmac.new(key = self.API_SECRET, digestmod = hashlib.sha512)
        H.update(data.encode('utf-8'))
        return H.hexdigest()

    def api_query(self, api_method, http_method = "POST", params = {}):
        params['nonce'] = int(round(time.time() * 1000))
        params =  urllib.parse.urlencode(params)
        
        sign = self.sha512(params)
        headers = {
            "Content-type": "application/x-www-form-urlencoded",
            "Key": self.API_KEY,
            "Sign": sign
        }
        conn = http.client.HTTPSConnection(self.API_URL)
        conn.request(http_method, "/" + self.API_VERSION + "/" + api_method, params, headers)
        response = conn.getresponse().read()
        
        conn.close()

        try:
            obj = json.loads(response.decode('utf-8'))
            if 'error' in obj and obj['error']:
                raise ScriptError(obj['error'])
            return obj
        except json.decoder.JSONDecodeError:
            raise ScriptError('An error occurred, while parsing response string:', response)

# Example
ExmoAPI_instance = ExmoAPI('YOUR API KEY', 'YOUR API SECRET')
print(ExmoAPI_instance.api_query('user_info'))
