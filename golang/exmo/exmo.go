/*
  Package exmo simplifies interaction with exmo.com API.
*/

package exmo

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// ApiResponse is a map for API responses
type ApiResponse map[string]interface{}

// ApiParams is a map for API calls' params
type ApiParams map[string]string

// Exmo holds client-specific info
type Exmo struct {
	key    string // public key
	secret string // secret key
}

// CheckErr checks error object
func CheckErr(err error) {
	if err != nil {
		log.Printf("api error: %s\n", err.Error())
	}
}

// Api creates Exmo instance with specified credentials
func Api(key string, secret string) Exmo {
	return Exmo{key, secret}
}

// Api_query is a general query method for API calls
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

// nonce generates request parameter ‘nonce’ with incremental numerical value (>0). The incremental numerical value should never reiterate or decrease.
func nonce() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// Do_sign encrypts POST data (param=val&param1=val1) with method HMAC-SHA512 using secret key; the secret key also can be found in user’s profile settings
func (ex *Exmo) Do_sign(message string) string {
	mac := hmac.New(sha512.New, []byte(ex.secret))
	mac.Write([]byte(message))
	return fmt.Sprintf("%x", mac.Sum(nil))
}

/*
   Public API
*/

// GetTrades return list of the deals on currency pairs
func (ex *Exmo) GetTrades(pair string) (response ApiResponse, err error) {
	response, err = ex.Api_query("public", "trades", ApiParams{"pair": pair})
	CheckErr(err)
	return
}

// GetOrderBook return the book of current orders on the currency pair
func (ex *Exmo) GetOrderBook(pair string, limit int) (response ApiResponse, err error) {
	if limit < 100 || limit > 1000 {
		fmt.Printf("limit param must be in range of 100-1000")
		response = nil
		err = errors.New("limit param must be in range of 100-1000")

	} else {
		response, err = ex.Api_query("public", "order_book", ApiParams{"pair": pair, "limit": string(limit)})
		CheckErr(err)
	}
	return
}

// Ticker return statistics on prices and volume of trades by currency pairs
func (ex *Exmo) Ticker() (response ApiResponse, err error) {
	response, err = ex.Api_query("public", "ticker", ApiParams{})
	if err != nil {
		fmt.Printf("api error: %s\n", err.Error())
	}
	return
}

// GetPairSettings return currency pairs settings
func (ex *Exmo) GetPairSettings() (response ApiResponse, err error) {
	response, err = ex.Api_query("public", "pair_settings", ApiParams{})
	CheckErr(err)
	return
}

// GetCurrency return currencies list
func (ex *Exmo) GetCurrency() (response []string, err error) {

	resp, err := http.Get("https://api.exmo.com/v1/currency")
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

	var dat []string
	err2 := json.Unmarshal([]byte(body), &dat)
	if err2 != nil {
		return nil, err2
	}

	return dat, nil
}

/*
   Authenticated API
*/

// GetUserInfo return information about user's account
func (ex *Exmo) GetUserInfo() (response ApiResponse, err error) {
	response, err = ex.Api_query("authenticated", "user_info", nil)
	CheckErr(err)
	return
}

// GetUserTrades return the list of user’s deals
func (ex *Exmo) GetUserTrades(pair string) (response ApiResponse, err error) {
	response, err = ex.Api_query("authenticated", "user_trades", ApiParams{"pair": pair})
	CheckErr(err)
	return
}

// OrderCreate creates order
func (ex *Exmo) OrderCreate(pair string, quantity string, price string, typeOrder string) (response ApiResponse, err error) {
	response, err = ex.Api_query("authenticated", "order_create", ApiParams{"pair": pair, "quantity": quantity, "price": price, "type": typeOrder})
	CheckErr(err)
	return
}

// Buy creates buy order
func (ex *Exmo) Buy(pair string, quantity string, price string) (response ApiResponse, err error) {
	response, err = ex.OrderCreate(pair, quantity, price, "buy")
	CheckErr(err)
	return
}

// Buy creates sell order
func (ex *Exmo) Sell(pair string, quantity string, price string) (response ApiResponse, err error) {
	response, err = ex.OrderCreate(pair, quantity, price, "sell")
	CheckErr(err)
	return
}

// MarketBuy creates market buy-order
func (ex *Exmo) MarketBuy(pair string, quantity string) (response ApiResponse, err error) {
	response, err = ex.OrderCreate(pair, quantity, "0", "market_buy")
	CheckErr(err)
	return
}

// MarketBuyTotal creates market buy-order for a certain amount (quantity parameter)
func (ex *Exmo) MarketBuyTotal(pair string, quantity string) (response ApiResponse, err error) {
	response, err = ex.OrderCreate(pair, quantity, "0", "market_buy_total")
	CheckErr(err)
	return
}

// MarketSell creates market sell-order
func (ex *Exmo) MarketSell(pair string, quantity string) (response ApiResponse, err error) {
	response, err = ex.OrderCreate(pair, quantity, "0", "market_sell")
	CheckErr(err)
	return
}

// MarketSellTotal creates market sell-order for a certain amount (quantity parameter)
func (ex *Exmo) MarketSellTotal(pair string, quantity string) (response ApiResponse, err error) {
	response, err = ex.OrderCreate(pair, quantity, "0", "market_sell_total")
	CheckErr(err)
	return
}

// OrderCancel cancels order
func (ex *Exmo) OrderCancel(orderId string) (response ApiResponse, err error) {
	response, err = ex.Api_query("authenticated", "order_cancel", ApiParams{"order_id": orderId})
	CheckErr(err)
	return
}

// GetUserOpenOrders returns the list of user’s active orders
func (ex *Exmo) GetUserOpenOrders() (response ApiResponse, err error) {
	response, err = ex.Api_query("authenticated", "user_open_orders", ApiParams{})
	CheckErr(err)
	return
}

// GetUserCancelledOrders returns the list of user’s deals
func (ex *Exmo) GetUserCancelledOrders(offset uint, limit uint) (response []interface{}, err error) {
	if limit < 100 || limit > 1000 {
		fmt.Printf("limit param must be in range of 100-1000")
		response = nil
		err = errors.New("limit param must be in range of 100-1000")
	} else {

		post_params := url.Values{}

		post_params.Add("nonce", nonce())
		post_params.Add("offset", string(offset))
		post_params.Add("limit", string(limit))

		post_content := post_params.Encode()

		sign := ex.Do_sign(post_content)

		req, _ := http.NewRequest("POST", "https://api.exmo.com/v1/user_cancelled_orders", bytes.NewBuffer([]byte(post_content)))
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

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(body), &response)
		if err != nil {
			return nil, err
		}
	}
	return
}

// GetOrderTrades returns the list of user’s cancelled orders
func (ex *Exmo) GetOrderTrades(orderId string) (response ApiResponse, err error) {
	response, err = ex.Api_query("authenticated", "order_trades", ApiParams{"order_id": orderId})
	CheckErr(err)
	return
}

// GetRequiredAmount calculating and returns the sum of buying a certain amount of currency for the particular currency pair
func (ex *Exmo) GetRequiredAmount(pair string, quantity string) (response ApiResponse, err error) {
	response, err = ex.Api_query("authenticated", "required_amount", ApiParams{"pair": pair, "quantity": quantity})
	CheckErr(err)
	return
}

// GetDepositAddress returns the list of addresses for cryptocurrency deposit
func (ex *Exmo) GetDepositAddress() (response ApiResponse, err error) {
	response, err = ex.Api_query("authenticated", "deposit_address", ApiParams{})
	CheckErr(err)
	return
}

/*
   Wallet API
*/

// GetWalletHistory returns history of wallet
func (ex *Exmo) GetWalletHistory(date time.Time) (response ApiResponse, err error) {
	dateUnix := date.Unix()

	dateConverted := strconv.Itoa(int(dateUnix))
	fmt.Println(dateConverted)

	if !date.IsZero() {
		response, err = ex.Api_query("authenticated", "wallet_history", ApiParams{})
	} else {
		response, err = ex.Api_query("authenticated", "wallet_history", ApiParams{"date": dateConverted})
	}

	CheckErr(err)
	return
}
