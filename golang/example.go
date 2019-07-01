package main

import (
	"fmt"
	"github.com/exmo-dev/exmo_api_lib/tree/master/golang/exmo"
	"math/big"
	"strconv"
	"time"
)

func main() {

	var orderId string

	key := ""    // TODO replace with your api key from profile page
	secret := "" // TODO replace with your api secret from profile page

	api := exmo.Api(key, secret)

	resultTrades, errTrades := api.GetTrades("BTC_RUB")
	if errTrades != nil {
		fmt.Errorf("api error: %s\n", errTrades.Error())
	} else {
		for _, v := range resultTrades {
			for k, val := range v.([]interface{}) {
				tmpindex := 0
				for key, value := range val.(map[string]interface{}) {
					if tmpindex != k {
						fmt.Printf("\n\nindex: %d \n", k)
						tmpindex = k
					}
					if key == "trade_id" {
						fmt.Println(key, big.NewFloat(value.(float64)).String())
					} else if key == "date" {
						fmt.Println(key, time.Unix(int64(value.(float64)), 0))
					} else {
						fmt.Println(key, value)
					}
				}
			}
		}
	}

	resultBook, errBook := api.GetOrderBook("BTC_RUB", 200)
	if errBook != nil {
		fmt.Errorf("api error: %s\n", errBook.Error())
	} else {
		for _, v := range resultBook {
			for key, value := range v.(map[string]interface{}) {
				if key == "bid" || key == "ask" {
					for _, val := range value.([]interface{}) {
						fmt.Printf("%s: ", key)
						for index, valnested := range val.([]interface{}) {
							switch index {
							case 0:
								fmt.Printf("price %s, ", valnested.(string))

							case 1:
								fmt.Printf("quantity %s, ", valnested.(string))
							case 2:
								fmt.Printf("total %s \n", valnested.(string))
							}
						}
					}
				} else {
					fmt.Println(key, value)
				}
			}

		}
	}

	ticker, errTicker := api.Ticker()
	if errTicker != nil {
		fmt.Printf("api error: %s\n", errTicker.Error())
	} else {
		for pair, pairvalue := range ticker {
			fmt.Printf("\n\n%s:\n", pair)
			for key, value := range pairvalue.(map[string]interface{}) {
				fmt.Println(key, value)
			}
		}
	}

	resultPairSettings, errPairSettings := api.GetPairSettings()
	if errPairSettings != nil {
		fmt.Printf("api error: %s\n", errPairSettings.Error())
	} else {
		for pair, pairvalue := range resultPairSettings {
			fmt.Printf("\n\n%s:\n", pair)
			for key, value := range pairvalue.(map[string]interface{}) {
				fmt.Println(key, value)
			}
		}
	}

	resultCurrency, errCurrency := api.GetCurrency()
	if errCurrency != nil {
		fmt.Printf("api error: %s\n", errCurrency.Error())
	} else {
		fmt.Println("\nCurrencies:")
		for _, pair := range resultCurrency {
			fmt.Println(pair)
		}
	}

	resultUserInfo, errUserInfo := api.GetUserInfo()
	if errUserInfo != nil {
		fmt.Printf("api error: %s\n", errUserInfo.Error())
	} else {
		for key, value := range resultUserInfo {
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

	usertrades, err1 := api.GetUserTrades("BTC_RUB")
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

	order, errOrder := api.Buy("BTC_RUB", "0.001", "50096")
	if errOrder != nil {
		fmt.Printf("api error: %s\n", errOrder.Error())
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
				val := strconv.Itoa(int(value.(float64)))
				orderId = val
				fmt.Printf("Order id: %s", orderId)
			}
		}
	}

	marketOrder, errMarketOrder := api.MarketBuy("BTC_RUB", "0.001")
	if errMarketOrder != nil {
		fmt.Printf("api error: %s\n", errMarketOrder.Error())
	} else {
		fmt.Println("Creating order...")
		for key, value := range marketOrder {
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

	orderSell, errOrderSell := api.Sell("BTC_RUB", "0.001", "800000")
	if errOrderSell != nil {
		fmt.Printf("api error: %s\n", errOrderSell.Error())
	} else {
		fmt.Println("Creating order...")
		for key, value := range orderSell {
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

	orderSellMarket, errOrderSellMarket := api.MarketSell("BTC_RUB", "0.001")
	if errOrderSellMarket != nil {
		fmt.Printf("api error: %s\n", errOrderSellMarket.Error())
	} else {
		fmt.Println("Creating order...")
		for key, value := range orderSellMarket {
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

	orderCancel, errCancel := api.OrderCancel(orderId)
	if errCancel != nil {
		fmt.Printf("api error: %s\n", errCancel.Error())
	} else {
		fmt.Printf("\nCancel order %s \n", orderId)
		for key, value := range orderCancel {
			if key == "result" && value != true {
				fmt.Println("\nError")
			}
			if key == "error" && value != "" {
				fmt.Println(value)
			}
		}
	}
}
