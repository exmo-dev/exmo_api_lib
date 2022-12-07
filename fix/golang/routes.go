package main

import (
	"log"

	"github.com/mikhalytch/eggs/math"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/fix44/executionreport"
	"github.com/quickfixgo/quickfix/fix44/marketdatarequestreject"
	"github.com/quickfixgo/quickfix/fix44/marketdatasnapshotfullrefresh"
	"github.com/quickfixgo/quickfix/fix44/securitydefinition"
	"github.com/quickfixgo/quickfix/fix44/securitylist"
)

func setupRoutes(router *quickfix.MessageRouter) {
	router.AddRoute(securitylist.Route(onSecurityList))
	router.AddRoute(securitydefinition.Route(onSecurityDefinition))
	router.AddRoute(marketdatasnapshotfullrefresh.Route(onMarketData))
	router.AddRoute(marketdatarequestreject.Route(onMarketDataReject))
	router.AddRoute(executionreport.Route(onExecutionReport))
}

func onExecutionReport(msg executionreport.ExecutionReport, _ quickfix.SessionID) quickfix.MessageRejectError {
	log.Println("execution report received:", limitMsg(msg))

	return nil
}

func onMarketDataReject(
	msg marketdatarequestreject.MarketDataRequestReject,
	_ quickfix.SessionID,
) quickfix.MessageRejectError {
	log.Println("market data reject received:", limitMsg(msg))

	return nil
}

func onMarketData(
	msg marketdatasnapshotfullrefresh.MarketDataSnapshotFullRefresh,
	_ quickfix.SessionID,
) quickfix.MessageRejectError {
	log.Println("market data received:", limitMsg(msg))

	return nil
}

func onSecurityDefinition(msg securitydefinition.SecurityDefinition, _ quickfix.SessionID) quickfix.MessageRejectError {
	log.Println("security definition received:", limitMsg(msg))

	return nil
}

func onSecurityList(msg securitylist.SecurityList, _ quickfix.SessionID) quickfix.MessageRejectError {
	log.Println("security list received:", limitMsg(msg))

	return nil
}

// -----

func limitMsg(msg quickfix.Messagable) string {
	const lengthLimit = 20

	str := msg.ToMessage().String()

	return str[:math.Min(len(str), lengthLimit)] + "..."
}
