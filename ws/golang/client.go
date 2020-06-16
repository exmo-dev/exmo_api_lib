package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	publicApiUsageExample()
}

func publicApiUsageExample() {
	startExmoClient("wss://ws-api.exmo.com:443/v1/public", []string{`{"id":1,"method":"subscribe","topics":["spot/trades:BTC_USD","spot/ticker:LTC_USD"]}`})
}

func privateApiUsageExample() {
	apiKey := ""
	secretKey := ""
	nonce := time.Now().In(time.UTC).Unix()

	mac := hmac.New(sha512.New, []byte(secretKey))
	_, _ = mac.Write([]byte(apiKey + strconv.Itoa(int(nonce))))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	loginCommand := fmt.Sprintf(`{"id":1,"method":"login","api_key":"%s","sign":"%s","nonce":%d}`, apiKey, sign, nonce)

	startExmoClient("wss://ws-api.exmo.com:443/v1/private", []string{loginCommand, `{"id":2,"method":"subscribe","topics":["spot/orders","spot/user_trades"]}`})
}

func startExmoClient(url string, initMessages []string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	for _, initMessage := range initMessages {
		err = c.WriteMessage(websocket.TextMessage, []byte(initMessage))
		if err != nil {
			log.Printf("fail to write init message: %v", err)
			return
		}
		log.Printf("sent: %s", initMessage)
	}

	select {
	case <-interrupt:
		log.Println("interrupt")

		// Cleanly close the connection by sending a close message and then
		// waiting (with timeout) for the server to close the connection.
		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("write close:", err)
			return
		}
		select {
		case <-done:
		case <-time.After(time.Second):
		}
	case <-done:
	}
}
