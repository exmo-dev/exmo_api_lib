package exmo

import (
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/exmo-dev/exmo_api_lib/tree/master/golang/exmo"
)

func TestApi_query(t *testing.T) {

	var orderId string // global variable for testing order cancelling after buying

	// ATTENTION!
	key := ""    // TODO replace with your api key from profile page
	secret := "" // TODO replace with your api secret from profile page

	api := exmo.Api(key, secret)

	t.Run("Get trades", func(t *testing.T) {
		result, err := api.GetTrades("BTC_RUB")
		if err != nil {
			fmt.Errorf("api error: %s\n", err.Error())
		} else {
			for _, v := range result {
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
	})

	t.Run("Get order book", func(t *testing.T) {
		result, err := api.GetOrderBook("BTC_RUB", 200)
		if err != nil {
			t.Errorf("api error: %s\n", err.Error())
		} else {
			for _, v := range result {

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
	})

	t.Run("Get ticker", func(t *testing.T) {
		ticker, errTicker := api.Ticker()
		if errTicker != nil {
			t.Errorf("api error: %s\n", errTicker.Error())
		} else {
			for pair, pairvalue := range ticker {
				fmt.Printf("\n\n%s:\n", pair)
				for key, value := range pairvalue.(map[string]interface{}) {
					fmt.Println(key, value)
				}
			}
		}
	})

	t.Run("Get pair settings", func(t *testing.T) {
		result, err := api.Ticker()
		if err != nil {
			t.Errorf("api error: %s\n", err.Error())
		} else {
			for pair, pairvalue := range result {
				fmt.Printf("\n\n%s:\n", pair)
				for key, value := range pairvalue.(map[string]interface{}) {
					fmt.Println(key, value)
				}
			}
		}
	})

	t.Run("Get currencies", func(t *testing.T) {
		result, err := api.GetCurrency()
		if err != nil {
			t.Errorf("api error: %s\n", err.Error())
		} else {
			fmt.Println("\nCurrencies:")
			for _, pair := range result {
				fmt.Println(pair)
			}
		}
	})

	t.Run("Get user info", func(t *testing.T) {
		fmt.Printf("-------------\n")
		result, err := api.GetUserInfo()
		if err != nil {
			t.Errorf("api error: %s\n", err.Error())
		} else {
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

	t.Run("Buy BTC (custom price)", func(t *testing.T) {
		order, errOrder := api.Buy("BTC_RUB", "0.001", "50096.72")
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
					val := strconv.Itoa(int(value.(float64)))
					orderId = val
					fmt.Printf("Order id: %s", orderId)
				}
			}
		}
	})

	t.Run("Buy BTC (market price)", func(t *testing.T) {
		order, errOrder := api.MarketBuy("BTC_RUB", "0.001")
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

	t.Run("Sell BTC (custom price)", func(t *testing.T) {
		order, errOrder := api.Sell("BTC_RUB", "0.001", "800000")
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

	t.Run("Sell BTC (market price)", func(t *testing.T) {
		order, errOrder := api.MarketSell("BTC_RUB", "0.001")
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

	t.Run("Cancel order", func(t *testing.T) {
		orderCancel, errCancel := api.OrderCancel(orderId)
		if errCancel != nil {
			t.Errorf("api error: %s\n", errCancel.Error())
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
	})

	t.Run("Get user's open orders", func(t *testing.T) {
		resultUserOpenOrders, errUserOpenOrders := api.GetUserOpenOrders()
		if errUserOpenOrders != nil {
			fmt.Errorf("api error: %s\n", errUserOpenOrders.Error())
		} else {
			for _, v := range resultUserOpenOrders {
				if v != nil {
					for _, val := range v.([]interface{}) {
						for key, value := range val.(map[string]interface{}) {
							if key == "quantity" {
								check, err := strconv.ParseFloat(value.(string), 64)
								if err != nil {
									t.Errorf("Could not convert %s to float64", key)
								}
								if check < 0 {
									t.Errorf("%s could not be less 0, got %d", key, value)
								}
							}
							if key == "price" {
								check, err := strconv.Atoi(value.(string))
								if err != nil {
									t.Errorf("Could not convert %s to int", key)
								}
								if check < 0 {
									t.Errorf("%s could not be less 0, got %d", key, value)
								}
							}
							if key == "amount" {
								check, err := strconv.ParseFloat(value.(string), 64)
								if err != nil {
									t.Errorf("Could not convert %s to float64", key)
								}
								if check < 0 {
									t.Errorf("%s could not be less 0, got %d", key, value)
								}
							}
						}
					}
				}
			}
		}
	})

	t.Run("Get user's cancelled orders", func(t *testing.T) {
		resultUserCancelledOrders, errUserCancelledOrders := api.GetUserCancelledOrders(0, 100)
		if errUserCancelledOrders != nil {
			fmt.Errorf("api error: %s\n", errUserCancelledOrders.Error())
		} else {
			for _, v := range resultUserCancelledOrders {

				if v != nil {
					for key, value := range v.(map[string]interface{}) {

						if key == "quantity" || key == "price" ||
							key == "amount" {
							check, ok := value.(float64)
							if ok != true {
								t.Errorf("Could not convert %s to float64", key)
							}
							if check < 0 {
								t.Errorf("%s could not be less 0, got %d", key, value)
							}
						}
					}
				}
			}
		}
	})

}
