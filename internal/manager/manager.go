package manager

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"runtime"
	"sync"
)

var gNum int

func init() {
	gNum = runtime.NumCPU()
}

type Option func(*Options)

func WithParallelNum(n int) Option {
	return func(m *Options) {
		m.gNum = n
	}
}

type Options struct {
	gNum int
}

func New(reader reader, logger logger, checker urlChecker, printer printer, o ...Option) *Manager {
	m := &Manager{reader: reader, logger: logger, fetcher: checker, printer: printer}
	m.opts.gNum = gNum

	for _, opt := range o {
		opt(&m.opts)
	}

	m.inputCh = make(chan url.URL, m.opts.gNum)

	return m
}

type Manager struct {
	opts    Options
	reader  reader
	logger  logger
	fetcher urlChecker
	printer printer
	inputCh chan url.URL
}

func (m *Manager) Run(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	m.makePool(ctx, wg)

	for {
		if err := ctx.Err(); err != nil {
			m.logger.Debug("context err", err.Error())
			return err
		}

		inputStr, err := m.reader.ReadURL()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			m.logger.Error("read url error: ", err.Error())
			continue
		}

		m.logger.Debug(fmt.Sprintf("read line url %s", inputStr))

		u, ok := ValidateURL(inputStr)
		if !ok {
			if _, err := m.printer.OutputValidationError(u.String()); err != nil {
				m.logger.Error("print error: ", err.Error())
				continue
			}
			m.logger.Debug(fmt.Sprintf("line with url %s is invalid", u.String()))
			continue
		}

		m.logger.Debug(fmt.Sprintf("url %s is corrected", u.String()))

		m.inputCh <- u
	}

	close(m.inputCh)

	wg.Wait()

	return nil
}

func (m *Manager) makePool(ctx context.Context, wg *sync.WaitGroup) {
	pNum := m.opts.gNum

	wg.Add(pNum)

	m.logger.Debug(fmt.Sprintf("start goroutine pool with %d goroutines", pNum))
	for i := 0; i < pNum; i++ {
		go m.poolHandler(ctx, wg)
	}
}

func (m *Manager) poolHandler(ctx context.Context, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	for ch := range m.inputCh {
		select {
		case <-ctx.Done():
			m.logger.Debug(ctx.Err().Error())
			return
		default:
		}

		m.logger.Debug(fmt.Sprintf("try check url %s", ch.String()))

		resp, err := m.fetcher.Check(ctx, ch.String())
		if err != nil {
			m.logger.Error("url check error: ", err.Error())
			continue
		}

		if _, err := m.printer.OutputEntry(
			ch.String(), resp.ContentLength, resp.StatusCode, resp.Status, resp.Ts,
		); err != nil {
			m.logger.Error("print error: ", err.Error())
			continue
		}
	}
}
