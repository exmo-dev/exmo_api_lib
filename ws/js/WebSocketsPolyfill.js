const WebSocketLib = require('ws');

function createWSConnection(props) {
  return typeof WebSocket !== 'undefined' ? new WebSocket(props) : new WebSocketLib(props);
}

module.exports = {
  createWSConnection,
};
