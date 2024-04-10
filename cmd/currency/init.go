package main

import (
	"context"
	"crypto/tls"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/crackc0der/currency/config"
	"github.com/crackc0der/currency/internal/currency"
	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
)

func Run() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	router := mux.NewRouter()
	scheduler := gocron.NewScheduler(time.UTC)
	timeout := 10
	idleTimeout := 15
	MaxHeaderBytes := 20
	readHeaderTimeout := 5

	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	dsn := config.GetDSN(conf)

	repository, err := currency.NewRepository(dsn)
	if err != nil {
		log.Fatal("error creating repository: ", err)
	}

	service := currency.NewService(repository, logger, conf)

	endpoint := currency.NewEndpoint(service, logger, conf)

	_, _ = scheduler.Every(conf.TimeOutUpdate).Seconds().Do(service.CurrencyMonitor)
	_, _ = scheduler.Every(conf.TimeOutUpdatePerHour).Seconds().Do(service.SetChangesPerHour)

	go scheduler.StartBlocking()

	router.HandleFunc("/rates", endpoint.GetCurrencies)
	router.HandleFunc("/rates/{name}", endpoint.GetCurrency)

	srv := http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    time.Duration(timeout) * time.Second,
		WriteTimeout:   time.Duration(timeout) * time.Second,
		IdleTimeout:    time.Duration(idleTimeout) * time.Second,
		MaxHeaderBytes: 1 << MaxHeaderBytes,
		ErrorLog:       log.New(os.Stderr, "http: ", log.LstdFlags),
		ConnState:      nil,
		TLSConfig:      nil,
		TLSNextProto:   make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		BaseContext: func(_ net.Listener) context.Context {
			return context.Background()
		},
		ConnContext: func(ctx context.Context, _ net.Conn) context.Context {
			return ctx
		},
		ReadHeaderTimeout:            time.Duration(readHeaderTimeout) * time.Second,
		DisableGeneralOptionsHandler: true,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
