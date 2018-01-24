package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/baskeboler/wordsoup-service/service"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/term"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)

	flag.Parse()

	var logger log.Logger
	{
		logger = term.NewColorLogger(os.Stderr, log.NewLogfmtLogger, func(keyvals ...interface{}) term.FgBgColor {
			for i := 0; i < len(keyvals)-1; i = i + 2 {
				if keyvals[i] == "err" && keyvals[i+1] != nil {
					return term.FgBgColor{Fg: term.Red}
				}
			}
			return term.FgBgColor{}
		})
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var s service.Service
	// {
	s, err := service.NewService()
	if err != nil {
		panic(err)
	}
	s = service.LoggingMiddleware(logger)(s)
	// }
	var h http.Handler
	{
		h = service.MakeHTTPRouter(s, log.With(logger, "component", "http"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {

		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("starting", "server")
	logger.Log("exit", <-errs)

}
