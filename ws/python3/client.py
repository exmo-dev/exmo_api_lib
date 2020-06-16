import asyncio
import base64
import hashlib
import hmac
import time

import websockets


def start_exmo_client(url, init_messages):
    async def ws_loop():
        async with websockets.connect(url) as websocket:
            for init_message in init_messages:
                await websocket.send(init_message)
                print("sent:", init_message)
            while True:
                print("read:", await websocket.recv())

    try:
        asyncio.get_event_loop().run_until_complete(ws_loop())
    except websockets.exceptions.ConnectionClosed as ex:
        print("connection closed", ex)
    except KeyboardInterrupt:
        pass


def public_api_usage_example():
    start_exmo_client(
        "wss://ws-api.exmo.com:443/v1/public",
        ["""{"id":1,"method":"subscribe","topics":["spot/trades:BTC_USD","spot/ticker:LTC_USD"]}"""],
    )


def private_api_usage_example():
    api_key = ""
    secret = ""
    nonce = str(int(time.time()))

    sign = hmac.new(secret.encode('utf8'), (api_key + nonce).encode('utf8'), hashlib.sha512).digest()
    sign = base64.b64encode(sign).decode('utf8')
    login_command = '{"id":1,"method":"login","api_key":"%s","sign":"%s","nonce":%s}' % (api_key, sign, nonce)

    start_exmo_client(
        "wss://ws-api.exmo.com:443/v1/private",
        [login_command, """{"id":2,"method":"subscribe","topics":["spot/orders","spot/user_trades"]}"""],
    )


if __name__ == "__main__":
    public_api_usage_example()
