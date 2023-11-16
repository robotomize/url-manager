package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/robotomize/url-manager/internal/httputil"
	"github.com/robotomize/url-manager/internal/manager"
	"github.com/robotomize/url-manager/internal/printer"
	"github.com/robotomize/url-manager/internal/urlchecker"
	"github.com/robotomize/url-manager/internal/urlreader"
)

var (
	source string
	debug  bool
	sync   bool
)

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&source,
		"source",
		"s",
		"",
		"source file with urls",
	)
	rootCmd.PersistentFlags().BoolVarP(
		&debug,
		"debug",
		"d",
		false,
		"debug logging",
	)
	rootCmd.PersistentFlags().BoolVarP(
		&sync,
		"sync",
		"c",
		false,
		"sync mode with one thread",
	)
}

var rootCmd = &cobra.Command{
	Use:          "url-manager",
	Long:         "Check websites console application",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer cancel()

		d := setup()

		l := d.logger
		if source == "" {
			l.Error("Source file with url is empty")
			return nil
		}

		f, err := os.OpenFile(source, os.O_RDONLY, 0x0660)
		if err != nil {
			switch {
			case errors.Is(err, os.ErrNotExist):
				l.Error("Source file not found")
			case errors.Is(err, os.ErrPermission):
				l.Error("Source file permission error")
			default:
				l.Error("Open file error:", err)
			}
			return nil
		}

		r := urlreader.New(f)

		opts := make([]manager.Option, 0)
		if sync {
			opts = append(opts, manager.WithParallelNum(1))
		}

		m := manager.New(r, l, d.urlChecker, d.printer, opts...)
		if err := m.Run(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				l.Error("process error:", err)
			}
		}
		return nil
	},
}

type deps struct {
	logger     *slog.Logger
	printer    printer.Printer
	httpClient httputil.Client
	urlChecker urlchecker.Checker
}

func setup() *deps {
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

	return d
}
