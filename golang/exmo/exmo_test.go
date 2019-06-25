package exmo

import (
	"fmt"
	"github.com/exmo-dev/exmo_api_lib/tree/master/golang/exmo"
	"testing"
)

func TestApi_query(t *testing.T) {

	// ATTENTION!
	key := ""    // TODO replace with your api key from profile page
	secret := "" // TODO replace with your api secret from profile page

	api := exmo.Api(key, secret)

	t.Run("Get user info", func(t *testing.T) {
		fmt.Printf("-------------\n")
		result, err := api.GetUserInfo()
		if err != nil {
			t.Errorf("api error: %s\n", err.Error())
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

		usertrades, err1 := api.GetUserTrades("BTC_RUB")
		if err1 != nil {
			t.Errorf("api error: %s\n", err1.Error())
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

	t.Run("Buy BTC", func(t *testing.T) {
		order, errOrder := api.Buy("BTC_RUB", "0.001", "664096.72")
		if errOrder != nil {
			t.Errorf("api error: %s\n", errOrder.Error())
		} else {
			fmt.Println("Creating order...")
			for key, value := range order {
				if key == "result" && value != true {
					fmt.Println("\nError")
				}
				if key == "error" && value != "" {
					fmt.Println(value)
				}
				if key == "order_id" && value != nil {
					fmt.Printf("Order id: %f", value.(float64))
				}
			}
		}
	})
}
