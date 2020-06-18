const CryptoJS = require('crypto-js');
const { createWSConnection } = require('./WebSocketsPolyfill');
const EXMO_WS_BASE_URL = `wss://ws-api.exmo.com:443/v1`;
const EXMO_WS_PUBLIC_URL = `${EXMO_WS_BASE_URL}/public`;
const EXMO_WS_PRIVATE_URL = `${EXMO_WS_BASE_URL}/private`;

function createExmoWSConnection(url, messages) {
  const socket = createWSConnection(url);
  const onMessage = (event) => console.log('message:', event);
  const onClose = (event) => console.log('connection closed:', event);
  const onError = (error) => console.log('connection error:', error);
  const onInitialize = () => {
    console.log('connection opened');

    for (let message of messages) {
      console.log('sending:', message);
      socket.send(message);
    }
  };

  socket.on('open', onInitialize);
  socket.on("message", onMessage);
  socket.on('close', onClose);
  socket.on('error', onError);
}

function connectExmoWSPublicApi() {
  const data = [
    '{"id":1,"method":"subscribe","topics":["spot/trades:BTC_USD","spot/ticker:LTC_USD"]}',
  ];

  createExmoWSConnection(EXMO_WS_PUBLIC_URL, data);
}

function connectExmoWSPrivateApi(apiKey) {
  const secretKey = '';
  const nonce = Date.now();
  const sign = CryptoJS.HmacSHA512(apiKey + nonce, secretKey).toString(CryptoJS.enc.Base64);
  const data = [
    `{"id":1,"method":"login","api_key":"${apiKey}","sign":"${sign}","nonce":${nonce}}`,
    '{"id":2,"method":"subscribe","topics":["spot/orders","spot/user_trades"]}',
  ];

  createExmoWSConnection(EXMO_WS_PRIVATE_URL, data);
}

module.exports = {
  connectExmoWSPublicApi,
  connectExmoWSPrivateApi,
}
