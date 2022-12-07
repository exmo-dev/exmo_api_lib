package main

import (
	"log"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/enum"
	"github.com/quickfixgo/quickfix/field"
	"github.com/quickfixgo/quickfix/fix44/marketdatarequest"
	"github.com/quickfixgo/quickfix/fix44/newordersingle"
	"github.com/quickfixgo/quickfix/fix44/securitydefinitionrequest"
	"github.com/quickfixgo/quickfix/fix44/securitylistrequest"
	"github.com/shopspring/decimal"
)

type Sender func(m quickfix.Messagable) error

// scenario performs imaginary actions sequence.
func scenario(session quickfix.SessionID) error {
	const pair = "BTC_USD"

	log.Println("scenario started")

	sender := func(m quickfix.Messagable) error {
		err := quickfix.SendToTarget(m, session)
		if err == nil {
			time.Sleep(time.Second) // since the scenario is way no interactive, and due to async implementation nature
		}

		return errors.Wrap(err, "unable to send quickfix msg")
	}

	if err := requestSecurityList(sender); err != nil {
		return errors.Wrap(err, "security list request")
	}

	if err := requestSecurityDefinition(sender, pair); err != nil {
		return errors.Wrap(err, "security definition request")
	}

	mdRequestID := strconv.Itoa(int(time.Now().UTC().Unix()))
	if err := requestMarketData(sender, mdRequestID, pair); err != nil {
		return errors.Wrap(err, "market data request")
	}

	if err := createBuyOrder(sender, pair); err != nil {
		return errors.Wrap(err, "create buy order")
	}

	if err := unsubscribeFromMarketData(sender, mdRequestID); err != nil {
		return errors.Wrap(err, "market data unsubscription")
	}

	if err := createSellOrder(sender, pair); err != nil {
		return errors.Wrap(err, "create sell order")
	}

	log.Println("scenario finished")

	return nil
}

func unsubscribeFromMarketData(sender Sender, mdRequestID string) error {
	log.Println("market data unsubscription")

	return sender(marketdatarequest.New(
		field.NewMDReqID(mdRequestID),
		field.NewSubscriptionRequestType(enum.SubscriptionRequestType_DISABLE_PREVIOUS_SNAPSHOT_PLUS_UPDATE_REQUEST),
		field.NewMarketDepth(0), // ignored
	))
}

func createBuyOrder(sender Sender, pair string) error {
	log.Println("creating buy order")

	return createOrder(sender, pair, enum.Side_BUY)
}

func createSellOrder(sender Sender, pair string) error {
	log.Println("creating sell order")

	return createOrder(sender, pair, enum.Side_SELL)
}

func createOrder(sender Sender, pair string, side enum.Side) error {
	const (
		quantityScale = 8
		priceScale    = 8
		price         = 15899
	)

	orderID := strconv.Itoa(int(time.Now().UTC().Unix()))

	singleOrderRequest := newordersingle.New(
		field.NewClOrdID(orderID),
		field.NewSide(side),
		field.NewTransactTime(time.Now().UTC()),
		field.NewOrdType(enum.OrdType_LIMIT),
	)
	singleOrderRequest.SetSymbol(pair)
	singleOrderRequest.SetOrderQty(decimal.New(1, -2), quantityScale) // 0.01
	singleOrderRequest.SetTimeInForce(enum.TimeInForce_GOOD_TILL_CANCEL)
	singleOrderRequest.SetExecInst("")
	singleOrderRequest.SetPrice(decimal.NewFromInt(price), priceScale)

	return sender(singleOrderRequest)
}

func requestMarketData(sender Sender, requestID string, pair string) error {
	log.Println("requesting market data")

	mdReq := marketdatarequest.New(
		field.NewMDReqID(requestID),
		field.NewSubscriptionRequestType(enum.SubscriptionRequestType_SNAPSHOT_PLUS_UPDATES),
		field.NewMarketDepth(0), // ignored
	)
	mdEntryTypesGrp := marketdatarequest.NewNoMDEntryTypesRepeatingGroup()
	mdEntryTypes := mdEntryTypesGrp.Add()
	mdEntryTypes.SetMDEntryType(enum.MDEntryType_BID)
	mdEntryTypes = mdEntryTypesGrp.Add()
	mdEntryTypes.SetMDEntryType(enum.MDEntryType_OFFER)
	mdReq.SetNoMDEntryTypes(mdEntryTypesGrp)

	relatedSymGrp := marketdatarequest.NewNoRelatedSymRepeatingGroup()
	relatedSym := relatedSymGrp.Add()
	relatedSym.SetSymbol(pair)
	mdReq.SetNoRelatedSym(relatedSymGrp)

	return sender(mdReq)
}

func requestSecurityDefinition(sender Sender, pair string) error {
	log.Println("requesting security definition")

	requestID := strconv.Itoa(int(time.Now().UTC().Unix()))
	securityDefinitionRequest := securitydefinitionrequest.New(
		field.NewSecurityReqID(requestID),
		field.NewSecurityRequestType(enum.SecurityRequestType_REQUEST_LIST_SECURITIES),
	)
	securityDefinitionRequest.SetSymbol(pair)

	return sender(securityDefinitionRequest)
}

func requestSecurityList(sender Sender) error {
	log.Println("requesting security list")

	requestID := strconv.Itoa(int(time.Now().UTC().Unix()))

	return sender(securitylistrequest.New(
		field.NewSecurityReqID(requestID),
		field.NewSecurityListRequestType(enum.SecurityListRequestType_ALL_SECURITIES)),
	)
}
