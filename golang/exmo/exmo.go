package exmo

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type ApiResponse map[string]interface{}
type ApiParams map[string]string

type Exmo struct {
	key    string
	secret string
}

func Api(key string, secret string) Exmo {
	return Exmo{key, secret}
}

func (ex *Exmo) Api_query(mode string, method string, params ApiParams) (ApiResponse, error) {

	post_params := url.Values{}
	if mode == "authenticated" {
		post_params.Add("nonce", nonce())
	}
	if params != nil {
		for key, value := range params {
			post_params.Add(key, value)
		}
	}
	post_content := post_params.Encode()

	sign := ex.Do_sign(post_content)

	req, _ := http.NewRequest("POST", "https://api.exmo.com/v1/"+method, bytes.NewBuffer([]byte(post_content)))
	req.Header.Set("Key", ex.key)
	req.Header.Set("Sign", sign)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(post_content)))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return nil, errors.New("http status: " + resp.Status)
	}

	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		return nil, err1
	}

	var dat map[string]interface{}
	err2 := json.Unmarshal([]byte(body), &dat)
	if err2 != nil {
		return nil, err2
	}

	if result, ok := dat["result"]; ok && result.(bool) != true {
		return nil, errors.New(dat["error"].(string))
	}

	return dat, nil
}

func nonce() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func (ex *Exmo) Do_sign(message string) string {
	mac := hmac.New(sha512.New, []byte(ex.secret))
	mac.Write([]byte(message))
	return fmt.Sprintf("%x", mac.Sum(nil))
}

// Get all trades
func (ex *Exmo) GetTrades(pair string) (response ApiResponse, err error) {
	response, err = ex.Api_query("authenticated", "trades", ApiParams{"pair": pair})
	if err != nil {
		fmt.Printf("api error: %s\n", err.Error())
	}
	return
}

// Get info about user account
func (ex *Exmo) GetUserInfo() (response ApiResponse, err error) {
	response, err = ex.Api_query("authenticated", "user_info", nil)
	if err != nil {
		fmt.Printf("api error: %s\n", err.Error())
	}
	return
}

func (ex *Exmo) GetUserTrades(pair string) (response ApiResponse, err error) {
	response, err = ex.Api_query("authenticated", "user_trades", ApiParams{"pair": pair})
	if err != nil {
		fmt.Printf("api error: %s\n", err.Error())
	}
	return
}

func (ex *Exmo) OrderCreate(pair string, quantity string, price string, typeOrder string) (response ApiResponse, err error) {
	response, err = ex.Api_query("authenticated", "order_create", ApiParams{"pair": pair, "quantity": quantity, "price": price, "type": typeOrder})
	if err != nil {
		fmt.Printf("api error: %s\n", err.Error())
	}
	return
}

func (ex *Exmo) Buy(pair string, quantity string, price string) (response ApiResponse, err error) {
	response, err = ex.OrderCreate(pair, quantity, price, "buy")
	if err != nil {
		fmt.Printf("api error: %s\n", err.Error())
	}
	return
}

func (ex *Exmo) Sell(pair string, quantity string, price string) (response ApiResponse, err error) {
	response, err = ex.OrderCreate(pair, quantity, price, "sell")
	if err != nil {
		fmt.Printf("api error: %s\n", err.Error())
	}
	return
}

func (ex *Exmo) MarketBuy(pair string, quantity string) (response ApiResponse, err error) {
	response, err = ex.OrderCreate(pair, quantity, "0", "market_buy")
	if err != nil {
		fmt.Printf("api error: %s\n", err.Error())
	}
	return
}

func (ex *Exmo) MarketSell(pair string, quantity string) (response ApiResponse, err error) {
	response, err = ex.OrderCreate(pair, quantity, "0", "market_sell")
	if err != nil {
		fmt.Printf("api error: %s\n", err.Error())
	}
	return
}
