import base64
import hashlib
import hmac
import time

from websocket import create_connection
from websocket._exceptions import WebSocketConnectionClosedException


def start_exmo_client(url, init_messages):
    ws = create_connection(url, connection='Connection: Upgrade')

    try:
        for init_message in init_messages:
            ws.send(init_message)
            print "sent:", init_message
        while True:
            print "read:", ws.recv()

    except WebSocketConnectionClosedException as ex:
        print "connection closed", ex
    except KeyboardInterrupt:
        pass
    finally:
        ws.close()


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
    private_api_usage_example()