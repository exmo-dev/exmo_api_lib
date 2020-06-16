<?php

namespace Exmo\WebSocketApi;

use Closure;

class Client
{
    const EVENT_INFO = 'info';
    const EVENT_ERROR = 'error';
    const EVENT_UPDATE = 'update';
    const EVENT_SNAPSHOT = 'snapshot';
    const EVENT_SUBSCRIBED = 'subscribed';
    const EVENT_UNSUBSCRIBED = 'unsubscribed';

    /**
     * @var \Ratchet\Client\WebSocket
     */
    protected $connect;

    /**
     * @var int
     */
    protected $messageId = 0;

    /**
     * @var string
     */
    protected $sessionId;

    /**
     * ExmoWsApi constructor.
     * @param \Ratchet\Client\WebSocket $connect
     */
    public function __construct(\Ratchet\Client\WebSocket $connect)
    {
        $this->connect = $connect;
    }

    /**
     * @return string|null
     */
    public function getSessionId()
    {
        return $this->sessionId;
    }

    /**
     * @param string $apiKey
     * @param string $apiSecret
     * @param null   $nonce
     */
    public function login($apiKey, $apiSecret, $nonce = null)
    {
        $nonce = $nonce ?: time();
        $this->send([
            'method'  => 'login',
            'api_key' => $apiKey,
            'nonce'   => $nonce,
            'sign'    => $this->getSign($apiKey, $apiSecret, $nonce),
        ]);
    }

    /**
     * @param array $topics
     */
    public function subscribe(array $topics)
    {
        $this->send([
            'method' => 'subscribe',
            'topics' => $topics,
        ]);
    }


    /**
     * @param array $topics
     */
    public function unsubscribe(array $topics)
    {
        $this->send([
            'method' => 'unsubscribe',
            'topics' => $topics,
        ]);
    }

    /**
     * @param array $data
     */
    public function send(array $data)
    {
        $this->connect->send(json_encode(
            [
                'id' => ++$this->messageId,
            ] + $data
        ));
    }

    /**
     * @param Closure $receiveCallback
     */
    public function onMessage(\Closure $receiveCallback)
    {
        $this->connect->on('message', function ($msg) use ($receiveCallback) {
            $receivedData = json_decode($msg, true);
            if (!empty($receivedData['session_id'])) {
                $this->sessionId = $receivedData['session_id'];
            }
            $receiveCallback($receivedData);
        });
    }

    /**
     * @param string $apiKey
     * @param string $apiSecret
     * @param int    $nonce
     * @return string
     */
    protected function getSign($apiKey, $apiSecret, $nonce)
    {
        return base64_encode(hash_hmac('sha512', $apiKey . $nonce, $apiSecret, true));
    }
}