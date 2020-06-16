<?php

require __DIR__ . '/vendor/autoload.php';

const EXMO_API_KEY = 'api_key_change_me';
const EXMO_API_SECRET = 'api_key_secret_change_me';

\Ratchet\Client\connect('wss://ws-api.exmo.com/v1/private')->then(function ($conn) {

    $exmoApi = new \Exmo\WebSocketApi\Client($conn);
    $exmoApi->login(EXMO_API_KEY, EXMO_API_SECRET);
    $exmoApi->subscribe([
        "spot/wallet",
        "spot/orders",
        "spot/user_trades",
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