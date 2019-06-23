package lib

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
	"testing"
	"time"
)

func TestApi_query(t *testing.T) {

	// ATTENTION!
	key := ""    // TODO replace with your api key from profile page
	secret := "" // TODO replace with your api secret from profile page

	type ApiResponse map[string]interface{}
	type ApiParams map[string]string

	nonce := func() string {
		t.Helper()
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	Do_sign := func(message string, secret string) string {
		t.Helper()
		mac := hmac.New(sha512.New, []byte(secret))
		mac.Write([]byte(message))
		return fmt.Sprintf("%x", mac.Sum(nil))
	}

	Api_query := func(key string, secret string, method string, params ApiParams) (ApiResponse, error) {
		t.Helper()
		post_params := url.Values{}
		post_params.Add("nonce", nonce())
		if params != nil {
			for key, value := range params {
				post_params.Add(key, value)
			}
		}
		post_content := post_params.Encode()

		sign := Do_sign(post_content, secret)

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

	t.Run("Get user info", func(t *testing.T) {
		fmt.Printf("-------------\n")
		result, err := Api_query(key, secret, "user_info", nil)
		if err != nil {
			fmt.Printf("api error: %s\n", err.Error())
		} else {
			fmt.Println("api result:")
			for key, value := range result {
				if key == "balances" {
					fmt.Println("\n-- balances:")
					for k, v := range value.(map[string]interface{}) {
						fmt.Println(k, v)
					}
				}
				if key == "reserved" {
					fmt.Println("\n-- reserved:")
					for k, v := range value.(map[string]interface{}) {
						fmt.Println(k, v)
					}
				}
			}

		}
	})

	t.Run("Get user trades", func(t *testing.T) {
		fmt.Printf("-------------\n")

		usertrades, err1 := Api_query(key, secret, "user_trades", ApiParams{"pair": "BTC_RUB"})
		if err1 != nil {
			fmt.Printf("api error: %s\n", err1.Error())
		} else {
			fmt.Println("User trades")
			for pair, val := range usertrades {
				fmt.Printf("\n\n %s", pair)
				for _, interfacevalue := range val.([]interface{}) {
					fmt.Printf("\n\n***\n")
					for k, v := range interfacevalue.(map[string]interface{}) {
						fmt.Println(k, v)
					}
				}
			}
		}
	})
}
