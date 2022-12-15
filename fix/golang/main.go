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
)

func main() {
	if err := main0(); err != nil {
		log.Fatalln(err)
	}
}

func main0() error {
	sessionsCh := make(chan quickfix.SessionID)
	defer close(sessionsCh)

	scenarioFinishedCh := make(chan struct{})
	defer close(scenarioFinishedCh)

	clientCloser, err := startClient(sessionsCh)
	if err != nil {
		return err
	}

	defer clientCloser()

	ctx, ctxCloser := getNotifiedCtx()
	defer ctxCloser()

	go func() {
		defer func() { scenarioFinishedCh <- struct{}{} }()

		session, err := awaitLogon(ctx, sessionsCh)
		if err != nil {
			log.Println(err)

			return
		}

		err = scenario(session)
		if err != nil {
			log.Println(err)
		}
	}()

	select {
	case <-ctx.Done():
	case <-scenarioFinishedCh:
	}

	return nil
}

func startClient(logonCh chan quickfix.SessionID) (func(), error) {
	configContents, err := os.ReadFile(path.Join("config", "config.ini"))
	if err != nil {
		return nil, errors.Wrap(err, "unable to read config")
	}

	appSettings, err := quickfix.ParseSettings(bytes.NewReader(configContents))
	if err != nil {
		return nil, fmt.Errorf("error reading ini: %w", err)
	}

	sessions := appSettings.SessionSettings()
	apiKeys := make(map[APIKey]APISecret, len(sessions))

	for _, sessionSettings := range sessions {
		sci, errS := sessionSettings.Setting(config.SenderCompID)
		sec, errP := sessionSettings.Setting("Password")

		if errS == nil && errP == nil {
			apiKeys[APIKey(sci)] = APISecret(sec)
		}
	}

	app := NewQuickFixApp(apiKeys, logonCh)
	screenLogFactory := quickfix.NewScreenLogFactory()

	if err != nil {
		return nil, fmt.Errorf("error creating file log factory: %w", err)
	}

	initiator, err := quickfix.NewInitiator(app, quickfix.NewMemoryStoreFactory(), appSettings, screenLogFactory)
	if err != nil {
		return nil, fmt.Errorf("unable to create Initiator: %w", err)
	}

	err = initiator.Start()
	if err != nil {
		return nil, fmt.Errorf("unable to start Initiator: %w", err)
	}

	printConfig(bytes.NewReader(configContents))

	return initiator.Stop, nil
}

func printConfig(reader io.Reader) {
	scanner := bufio.NewScanner(reader)

	color.Set(color.Bold)
	log.Println("Started FIX initiator with config:")
	color.Unset()

	color.Set(color.FgHiMagenta)

	for scanner.Scan() {
		line := scanner.Text()
		log.Println(line)
	}

	color.Unset()
}

func getNotifiedCtx() (context.Context, func()) {
	// Listen to interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	return ctx, func() {
		stop()

		if errors.Is(ctx.Err(), context.Canceled) {
			log.Println("FIX initiator is stopped")

			return
		}

		log.Println(ctx.Err())
	}
}

func awaitLogon(ctx context.Context, sessionsCh chan quickfix.SessionID) (quickfix.SessionID, error) {
	log.Println("awaiting logon...")
	select {
	case <-ctx.Done():
		return quickfix.SessionID{}, errors.Wrap(ctx.Err(), "context was closed while awaiting for logon")
	case session, open := <-sessionsCh:
		if !open {
			return quickfix.SessionID{}, errors.New("no more sessions, closing")
		}

		log.Println("logged on:", session)

		return session, nil
	}
}
