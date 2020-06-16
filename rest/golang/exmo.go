package main

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

func api_query(method string, params ApiParams) (ApiResponse, error) {
	key := ""       // TODO replace with your api key from profile page
	secret := ""    // TODO replace with your api secret from profile page

	post_params := url.Values{}
	post_params.Add("nonce", nonce())
	if params != nil {
		for key, value := range params {
			post_params.Add(key, value)
		}
	}
	post_content := post_params.Encode()

	sign := do_sign(post_content, secret)

	req, _ := http.NewRequest("POST", "https://api.exmo.com/v1/"+method, bytes.NewBuffer([]byte(post_content)))
	req.Header.Set("Key", key)
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

func do_sign(message string, secret string) string {
	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write([]byte(message))
	return fmt.Sprintf("%x", mac.Sum(nil))
}

func main() {
	result, err := api_query("user_info", nil)
	if err != nil {
		fmt.Printf("api error: %s\n", err.Error())
	} else {
		fmt.Printf("api result: %v\n", result)
	}

	fmt.Printf("-------------\n")

	result1, err1 := api_query("user_trades", ApiParams{"pair": "BTC_USD"})
	if err1 != nil {
		fmt.Printf("api error: %s\n", err1.Error())
	} else {
		fmt.Printf("api result: %v\n", result1)
	}
}
