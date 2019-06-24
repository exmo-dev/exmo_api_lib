package main

import (
	"fmt"

	"github.com/exmo-dev/exmo_api_lib/tree/master/golang/lib"
)

func main() {

	key := ""    // TODO replace with your api key from profile page
	secret := "" // TODO replace with your api secret from profile page

	result, err := lib.GetUserInfo(key, secret)
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

	fmt.Printf("-------------\n")

	usertrades, err1 := lib.GetUserTrades(key, secret, lib.ApiParams{"pair": "BTC_RUB"})
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
}
