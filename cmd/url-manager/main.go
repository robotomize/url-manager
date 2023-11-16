package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/robotomize/url-manager/internal/httputil"
	"gitlab.com/robotomize/url-manager/internal/manager"
	"gitlab.com/robotomize/url-manager/internal/printer"
	"gitlab.com/robotomize/url-manager/internal/urlchecker"
	"gitlab.com/robotomize/url-manager/internal/urlreader"
)

var (
	source string
	debug  bool
	sync   bool
)

func init() {
	flag.StringVar(&source, "s", "", "-s filepath")
	flag.BoolVar(&debug, "d", false, "-d")
	flag.BoolVar(&sync, "sync", false, "-sync")
	flag.Parse()
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	d, err := setup()
	if err != nil {
		fmt.Printf("Unhandler error")
		os.Exit(1)
		return
	}

	l := d.logger
	if source == "" {
		l.Error("Source file with url is empty")
		return
	}

	f, err := os.OpenFile(source, os.O_RDONLY, 0x0660)
	if err != nil {
		switch {
		case errors.Is(err, os.ErrNotExist):
			l.Error("Source file not found")
		case errors.Is(err, os.ErrPermission):
			l.Error("Source file permission error")
		default:
			l.Error("Open file error: ", err.Error())
		}
		return
	}

	r := urlreader.New(f)
	var opts []manager.Option
	if sync {
		opts = append(opts, manager.WithParallelNum(1))
	}

	m := manager.New(r, l, d.urlChecker, d.printer)
	if err := m.Run(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			l.Error("process error: ", err.Error())
		}
	}
}

type deps struct {
	logger     *slog.Logger
	printer    printer.Printer
	httpClient httputil.Client
	urlChecker urlchecker.Checker
}

func setup() (*deps, error) {
	const (
		defaultClientTimeout = 10 * time.Second
		defaultRetryCount    = 3
		defaultRetryMinWait  = 2 * time.Second
		defaultRetryMaxWait  = 10 * time.Second
	)

	d := &deps{}

	logLevel := slog.LevelError
	if debug {
		logLevel = slog.LevelDebug
	}

	textLogger := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: logLevel,
			},
		),
	)

	stdoutPrinter := printer.New(os.Stdout)

	client := http.DefaultClient
	client.Timeout = defaultClientTimeout

	retryClient := httputil.NewRetryClient(
		http.DefaultClient, defaultRetryCount, defaultRetryMinWait, defaultRetryMaxWait, httputil.ExponentialBackoff,
	)

	checker := urlchecker.New(retryClient)

	d.urlChecker = checker
	d.httpClient = retryClient
	d.printer = stdoutPrinter
	d.logger = textLogger

	return d, nil
}
