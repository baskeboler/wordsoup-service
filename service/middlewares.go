package service

import (
	"context"
	"time"

	"github.com/baskeboler/wordsoup"
	"github.com/go-kit/kit/log"
)

// Middleware definition
type Middleware func(Service) Service

// LoggingMiddleware returns a middleware that logs each service call
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (m *loggingMiddleware) GeneratePuzzle(c context.Context, width, height, words int) (ws *wordsoup.WordSoup, err error) {
	defer func(begin time.Time) {
		m.logger.Log("method", "GeneratePuzzle", "took", time.Since(begin), "err", err)
	}(time.Now())
	ws, err = m.next.GeneratePuzzle(c, width, height, words)
	return
}
