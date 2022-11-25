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
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/config"
	"github.com/quickfixgo/quickfix/enum"
)

type (
	ApiKey    string
	ApiSecret string
)

// Client implements the quickfix.Application interface
type Client struct {
	keys map[ApiKey]ApiSecret
}

// OnCreate implemented as part of Application interface
func (e Client) OnCreate(sessionID quickfix.SessionID) {}

// OnLogon implemented as part of Application interface
func (e Client) OnLogon(sessionID quickfix.SessionID) {}

// OnLogout implemented as part of Application interface
func (e Client) OnLogout(sessionID quickfix.SessionID) {}

// FromAdmin implemented as part of Application interface
func (e Client) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	if msg.IsMsgTypeOf(enum.MsgType_LOGON) {
		log.Printf("api key [%s] is logged on", sessionID.SenderCompID)
	}
	return nil
}

// ToAdmin implemented as part of Application interface
func (e Client) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {
	apiKey := ApiKey(sessionID.SenderCompID)
	secret, exist := e.keys[apiKey]
	if !exist {
		log.Fatalf("unknown api-key [%s] in sessionID", apiKey)
	}
	SignLogonMsg(msg, secret)
}

// ToApp implemented as part of Application interface
func (e Client) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) (err error) {
	fmt.Printf("Sending %s\n", msg)
	return
}

// FromApp implemented as part of Application interface. This is the callback for all Application level messages from the counter party.
func (e Client) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	fmt.Printf("FromApp: %s\n", msg.String())
	return
}

func main() {
	err := startClient()
	if err != nil {
		log.Fatal(err)
	}
}

func startClient() error {
	cfgFileName := path.Join("config", "config.cfg")

	cfg, err := os.Open(cfgFileName)
	if err != nil {
		return fmt.Errorf("Error opening %v, %v\n", cfgFileName, err)
	}
	defer cfg.Close()

	stringData, readErr := io.ReadAll(cfg)
	if readErr != nil {
		return fmt.Errorf("error reading cfg: %s", readErr)
	}

	appSettings, err := quickfix.ParseSettings(bytes.NewReader(stringData))
	if err != nil {
		return fmt.Errorf("error reading cfg: %s", err)
	}

	sessions := appSettings.SessionSettings()
	apiKeys := make(map[ApiKey]ApiSecret, len(sessions))
	for _, sessionSettings := range sessions {
		sci, serr := sessionSettings.Setting(config.SenderCompID)
		sec, perr := sessionSettings.Setting("Password")
		if serr == nil && perr == nil {
			apiKeys[ApiKey(sci)] = ApiSecret(sec)
		}
	}
	app := Client{keys: apiKeys}
	screenLogFactory := quickfix.NewScreenLogFactory()

	if err != nil {
		return fmt.Errorf("error creating file log factory: %s", err)
	}

	initiator, err := quickfix.NewInitiator(app, quickfix.NewMemoryStoreFactory(), appSettings, screenLogFactory)
	if err != nil {
		return fmt.Errorf("Unable to create Initiator: %s\n", err)
	}

	err = initiator.Start()
	if err != nil {
		return fmt.Errorf("Unable to start Initiator: %s\n", err)
	}

	printConfig(bytes.NewReader(stringData))

	defer initiator.Stop()
	awaitTermination()
	return nil
}

func printConfig(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	color.Set(color.Bold)
	fmt.Println("Started FIX initiator with config:")
	color.Unset()

	color.Set(color.FgHiMagenta)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	color.Unset()
}

func awaitTermination() {
	// Listen to interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer func() {
		stop()
		if errors.Is(ctx.Err(), context.Canceled) {
			log.Println("FIX initiator is stopped")
			return
		}
		log.Println(ctx.Err())
	}()
	<-ctx.Done()
}
