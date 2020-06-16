<?php

require __DIR__ . '/vendor/autoload.php';

\Ratchet\Client\connect('wss://ws-api.exmo.com/v1/public')->then(function ($conn) {

    $exmoApi = new \Exmo\WebSocketApi\Client($conn);
    $exmoApi->subscribe([
        "spot/trades:BTC_USD",
        "spot/ticker:LTC_USD",
    ]);

    $exmoApi->onMessage(function ($data) {
        if ($data['event'] === \Exmo\WebSocketApi\Client::EVENT_ERROR) {
            throw new Exception($data['message'], $data['code']);
        }
        print_r($data);
    });

    $conn->on('close', function($code = null, $reason = null) use ($exmoApi) {
        echo "Connection closed. Code: $code; Reason: {$reason}; Session ID: {$exmoApi->getSessionId()}" . PHP_EOL;
    });

}, function (\Exception $e) {
    echo "Could not connect: {$e->getMessage()}" . PHP_EOL;
});