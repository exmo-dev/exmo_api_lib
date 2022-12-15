// Copyright (c) quickfixengine.org  All rights reserved.
//
// This file may be distributed under the terms of the quickfixengine.org
// license as defined by quickfixengine.org and appearing in the file
// LICENSE included in the packaging of this file.
//
// This file is provided AS IS with NO WARRANTY OF ANY KIND, INCLUDING
// THE WARRANTY OF DESIGN, MERCHANTABILITY AND FITNESS FOR A
// PARTICULAR PURPOSE.
//
// See http://www.quickfixengine.org/LICENSE for licensing information.
//
// Contact ask@quickfixengine.org if any conditions of this licensing
// are not clear to you.

package main

import (
	"log"

	"github.com/pkg/errors"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/enum"
)

type (
	APIKey    string
	APISecret string
)

// QuickFixApp implements the quickfix.Application interface.
type QuickFixApp struct {
	*quickfix.MessageRouter
	keys    map[APIKey]APISecret
	logonCh chan quickfix.SessionID
}

func NewQuickFixApp(keys map[APIKey]APISecret, logonSessionsCh chan quickfix.SessionID) *QuickFixApp {
	res := &QuickFixApp{
		MessageRouter: quickfix.NewMessageRouter(),
		keys:          keys,
		logonCh:       logonSessionsCh,
	}
	setupRoutes(res.MessageRouter)

	return res
}

// OnCreate implemented as part of Application interface.
func (e *QuickFixApp) OnCreate(quickfix.SessionID) {}

// OnLogon implemented as part of Application interface.
func (e *QuickFixApp) OnLogon(sessionID quickfix.SessionID) { e.logonCh <- sessionID }

// OnLogout implemented as part of Application interface.
func (e *QuickFixApp) OnLogout(quickfix.SessionID) {}

// FromAdmin implemented as part of Application interface.
func (e *QuickFixApp) FromAdmin(
	msg *quickfix.Message,
	sessionID quickfix.SessionID,
) quickfix.MessageRejectError {
	if msg.IsMsgTypeOf(enum.MsgType_LOGON) {
		log.Printf("api key [%s] is logged on", sessionID.SenderCompID)
	}

	return nil
}

// ToAdmin implemented as part of Application interface.
func (e *QuickFixApp) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {
	apiKey := APIKey(sessionID.SenderCompID)

	secret, exist := e.keys[apiKey]
	if !exist {
		log.Fatalf("unknown api-key [%s] in sessionID", apiKey)
	}

	SignLogonMsg(msg, secret)
}

// ToApp implemented as part of Application interface.
func (e *QuickFixApp) ToApp(msg *quickfix.Message, _ quickfix.SessionID) error {
	log.Printf("Sending %s\n", msg)

	return nil
}

// FromApp implemented as part of Application interface.
// This is the callback for all Application level messages from the counterparty.
func (e *QuickFixApp) FromApp(
	msg *quickfix.Message,
	sessionID quickfix.SessionID,
) quickfix.MessageRejectError {
	err := e.Route(msg, sessionID)
	if err != nil {
		if errors.Is(err, quickfix.UnsupportedMessageType()) {
			log.Printf("FromApp: %s\n", msg.String())
		} else {
			log.Fatalln(err)
		}
	}

	return nil
}
