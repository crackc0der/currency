package main

import (
	"exchange_course/config"
	"exchange_course/internal/currency"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

func Run() {
	configFile, err := os.ReadFile("../../config/config.yaml")
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	r := mux.NewRouter()

	if err != nil {
		log.Fatal("could not read config file: ", err)
	}

	config := config.Config{}
	err = yaml.Unmarshal(configFile, &config)

	if err != nil {
		log.Fatal("could not unmarshal config file: ", err)
	}

	dsn := "postgres://" + config.DataBase.DbUser + ":" + config.DataBase.DbPassword + "@" + config.DataBase.DbHost + ":" + config.DataBase.DbPort + "/" + config.DataBase.DbName + "?sslmode=disable"
	repository, err := currency.NewRepository(dsn)

	if err != nil {
		log.Fatal("error creating repository: ", err)
	}

	service := currency.NewService(repository)

	endpoint := currency.NewEndpoint(service, logger, &config)

	endpoint.CurrencyMonitor()

	r.HandleFunc("/", endpoint.GetCurrencies)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
