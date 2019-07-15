package exmo

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/exmo-dev/exmo_api_lib/tree/master/golang/exmo"
)

func TestApi_query(t *testing.T) {

	//var orderId string // global variable for testing order cancelling after buying

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
				for _, val := range v.([]interface{}) {
					for key, value := range val.(map[string]interface{}) {
						if key == "trade_id" || key == "date" {
							check, ok := value.(float64)
							if !ok {
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

	t.Run("Get order book", func(t *testing.T) {
		result, err := api.GetOrderBook("BTC_RUB", 200)
		if err != nil {
			t.Errorf("api error: %s\n", err.Error())
		} else {
			for _, v := range result {
				for key, value := range v.(map[string]interface{}) {
					if key == "bid" || key == "ask" {
						for _, val := range value.([]interface{}) {
							for _, valnested := range val.([]interface{}) {
								check, err := strconv.ParseFloat(valnested.(string), 64)
								if err != nil {
									t.Errorf("Could not convert %s to float64", key)
								}
								if check < 0 {
									t.Errorf("%s could not be less 0, got %d", key, valnested)
								}
							}
						}
					} else {
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
	})

	t.Run("Get ticker", func(t *testing.T) {
		ticker, errTicker := api.Ticker()
		if errTicker != nil {
			t.Errorf("api error: %s\n", errTicker.Error())
		} else {
			for _, pairvalue := range ticker {
				for key, value := range pairvalue.(map[string]interface{}) {
					if key == "updated" {
						check, ok := value.(float64)
						if !ok {
							t.Errorf("Could not convert %s to float64", key)
						}
						if check < 0 {
							t.Errorf("%s could not be less 0, got %d", key, value)
						}
					} else {
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
	})

	t.Run("Get pair settings", func(t *testing.T) {
		result, err := api.Ticker()
		if err != nil {
			t.Errorf("api error: %s\n", err.Error())
		} else {
			for _, pairvalue := range result {
				for key, value := range pairvalue.(map[string]interface{}) {
					if key == "updated" {
						check, ok := value.(float64)
						if !ok {
							t.Errorf("Could not convert %s to float64", key)
						}
						if check < 0 {
							t.Errorf("%s could not be less 0, got %d", key, value)
						}
					} else {
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
	})

	t.Run("Get currencies", func(t *testing.T) {
		result, err := api.GetCurrency()
		if err != nil {
			t.Errorf("api error: %s\n", err.Error())
		} else {
			for _, pair := range result {
				if reflect.TypeOf(pair).Name() != "string" {
					t.Errorf("response item %#v not a string", pair)
				}
			}
		}
	})

	t.Run("Get user info", func(t *testing.T) {
		result, err := api.GetUserInfo()
		if err != nil {
			t.Errorf("api error: %s\n", err.Error())
		} else {
			for key, value := range result {
				if key == "balances" || key == "reserved" {
					for k, v := range value.(map[string]interface{}) {
						check, err := strconv.ParseFloat(v.(string), 64)
						if err != nil {
							t.Errorf("Could not convert %s to float64", k)
						}
						if check < 0 {
							t.Errorf("%s could not be less 0, got %d", k, v)
						}
					}
				} else {
					check, ok := value.(float64)
					if !ok {
						t.Errorf("Could not convert %s to float64", key)
					}
					if check < 0 {
						t.Errorf("%s could not be less 0, got %d", key, value)
					}
				}
			}

		}
	})

	t.Run("Get user trades", func(t *testing.T) {
		usertrades, err1 := api.GetUserTrades("BTC_RUB")
		if err1 != nil {
			t.Errorf("api error: %s\n", err1.Error())
		} else {
			for _, val := range usertrades {
				for _, interfacevalue := range val.([]interface{}) {
					for k, v := range interfacevalue.(map[string]interface{}) {
						if k == "trade_id" || k == "date" || k == "order_id" {
							check, ok := v.(float64)
							if !ok {
								t.Errorf("Could not convert %s to float64", k)
							}
							if check < 0 {
								t.Errorf("%s could not be less 0, got %d", k, v)
							}
						} else if k == "quantity" || k == "price" || k == "amount" {
							check, err := strconv.ParseFloat(v.(string), 64)
							if err != nil {
								t.Errorf("Could not convert %s to float64", k)
							}
							if check < 0 {
								t.Errorf("%s could not be less 0, got %d", k, v)
							}
						} else {
							if reflect.TypeOf(v).Name() != "string" {
								t.Errorf("response item %s (value %#v) not a string, but %T", k, v, v)
							}
						}
					}
				}
			}
		}
	})

	//t.Run("Buy BTC (custom price)", func(t *testing.T) {
	//	order, errOrder := api.Buy("BTC_RUB", "0.001", "50096.72")
	//	if errOrder != nil {
	//		t.Errorf("api error: %s\n", errOrder.Error())
	//	} else {
	//		fmt.Println("Creating order...")
	//		for key, value := range order {
	//			if key == "result" && value != true {
	//				fmt.Println("\nError")
	//			}
	//			if key == "error" && value != "" {
	//				fmt.Println(value)
	//			}
	//			if key == "order_id" && value != nil {
	//				val := strconv.Itoa(int(value.(float64)))
	//				orderId = val
	//				fmt.Printf("Order id: %s", orderId)
	//			}
	//		}
	//	}
	//})
	//
	//t.Run("Buy BTC (market price)", func(t *testing.T) {
	//	order, errOrder := api.MarketBuy("BTC_RUB", "0.001")
	//	if errOrder != nil {
	//		t.Errorf("api error: %s\n", errOrder.Error())
	//	} else {
	//		fmt.Println("Creating order...")
	//		for key, value := range order {
	//			if key == "result" && value != true {
	//				fmt.Println("\nError")
	//			}
	//			if key == "error" && value != "" {
	//				fmt.Println(value)
	//			}
	//			if key == "order_id" && value != nil {
	//				fmt.Printf("Order id: %f", value.(float64))
	//			}
	//		}
	//	}
	//})
	//
	//t.Run("Sell BTC (custom price)", func(t *testing.T) {
	//	order, errOrder := api.Sell("BTC_RUB", "0.001", "800000")
	//	if errOrder != nil {
	//		t.Errorf("api error: %s\n", errOrder.Error())
	//	} else {
	//		fmt.Println("Creating order...")
	//		for key, value := range order {
	//			if key == "result" && value != true {
	//				fmt.Println("\nError")
	//			}
	//			if key == "error" && value != "" {
	//				fmt.Println(value)
	//			}
	//			if key == "order_id" && value != nil {
	//				fmt.Printf("Order id: %f", value.(float64))
	//			}
	//		}
	//	}
	//})
	//
	//t.Run("Sell BTC (market price)", func(t *testing.T) {
	//	order, errOrder := api.MarketSell("BTC_RUB", "0.001")
	//	if errOrder != nil {
	//		t.Errorf("api error: %s\n", errOrder.Error())
	//	} else {
	//		fmt.Println("Creating order...")
	//		for key, value := range order {
	//			if key == "result" && value != true {
	//				fmt.Println("\nError")
	//			}
	//			if key == "error" && value != "" {
	//				fmt.Println(value)
	//			}
	//			if key == "order_id" && value != nil {
	//				fmt.Printf("Order id: %f", value.(float64))
	//			}
	//		}
	//	}
	//})

	//t.Run("Cancel order", func(t *testing.T) {
	//	orderCancel, errCancel := api.OrderCancel(orderId)
	//	if errCancel != nil {
	//		t.Errorf("api error: %s\n", errCancel.Error())
	//	} else {
	//		fmt.Printf("\nCancel order %s \n", orderId)
	//		for key, value := range orderCancel {
	//			if key == "result" && value != true {
	//				fmt.Println("\nError")
	//			}
	//			if key == "error" && value != "" {
	//				fmt.Println(value)
	//			}
	//		}
	//	}
	//})
	//
	//t.Run("Get user's open orders", func(t *testing.T) {
	//	resultUserOpenOrders, errUserOpenOrders := api.GetUserOpenOrders()
	//	if errUserOpenOrders != nil {
	//		fmt.Errorf("api error: %s\n", errUserOpenOrders.Error())
	//	} else {
	//		for _, v := range resultUserOpenOrders {
	//			if v != nil {
	//				for _, val := range v.([]interface{}) {
	//					for key, value := range val.(map[string]interface{}) {
	//						if key == "quantity" {
	//							check, err := strconv.ParseFloat(value.(string), 64)
	//							if err != nil {
	//								t.Errorf("Could not convert %s to float64", key)
	//							}
	//							if check < 0 {
	//								t.Errorf("%s could not be less 0, got %d", key, value)
	//							}
	//						}
	//						if key == "price" {
	//							check, err := strconv.Atoi(value.(string))
	//							if err != nil {
	//								t.Errorf("Could not convert %s to int", key)
	//							}
	//							if check < 0 {
	//								t.Errorf("%s could not be less 0, got %d", key, value)
	//							}
	//						}
	//						if key == "amount" {
	//							check, err := strconv.ParseFloat(value.(string), 64)
	//							if err != nil {
	//								t.Errorf("Could not convert %s to float64", key)
	//							}
	//							if check < 0 {
	//								t.Errorf("%s could not be less 0, got %d", key, value)
	//							}
	//						}
	//					}
	//				}
	//			}
	//		}
	//	}
	//})

	t.Run("Get user's cancelled orders", func(t *testing.T) {
		resultUserCancelledOrders, errUserCancelledOrders := api.GetUserCancelledOrders(0, 100)
		if errUserCancelledOrders != nil {
			fmt.Errorf("api error: %s\n", errUserCancelledOrders.Error())
		} else {
			for _, v := range resultUserCancelledOrders {

				if v != nil {
					for key, value := range v.(map[string]interface{}) {

						if key == "quantity" || key == "price" || key == "amount" {
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

	t.Run("Calculating the sum of buying a certain amount of currency for the particular currency pair", func(t *testing.T) {
		resultRequiredAmount, errRequiredAmount := api.GetRequiredAmount("BTC_RUB", "0.01")
		if errRequiredAmount != nil {
			fmt.Errorf("api error: %s\n", errRequiredAmount.Error())
		} else {
			for k, v := range resultRequiredAmount {
				check, err := strconv.ParseFloat(v.(string), 64)
				if err != nil {
					t.Errorf("Could not convert %s to float64", k)
				}
				if check < 0 {
					t.Errorf("%s could not be less 0, got %d", k, v)
				}
			}
		}
	})

	t.Run("Getting the list of addresses for cryptocurrency deposit", func(t *testing.T) {
		resultDepositAddress, errDepositAddress := api.GetDepositAddress()
		if errDepositAddress != nil {
			fmt.Errorf("api error: %s\n", errDepositAddress.Error())
		} else {
			for _, v := range resultDepositAddress {
				_, ok := v.(string)
				if ok != true {
					t.Errorf("Could not convert %s address to string", key)
				}
			}
		}
	})
}
