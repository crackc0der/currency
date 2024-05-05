package currency

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"gopkg.in/telebot.v3"
)

func getCurrencies(service *Service) (string, error) {
	var message string
	currencies, err := service.GetCurrencies(context.Background())
	if err != nil {
		fmt.Printf("error in bot handle /rates: %v", err)
		return "", fmt.Errorf("error im method getCurrencies: %v", err)
	}

	for _, currency := range currencies {
		message += fmt.Sprintf("%s = %.2f ", currency.CurrencyName, currency.CurrencyPrice)
	}

	return message, nil
}

func BotRun(key string, service *Service) {
	autoChan := make(chan struct{})
	var timeout time.Duration

	pref := telebot.Settings{
		Token:  key,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/start", func(c telebot.Context) error {
		return c.Send(`The bot supports several commands. The /rates command without parameters will display the BTC and ETH rates. 
		The /rates command with the BTC or ETH parameter will display the rate of the selected currency. 
		The /start_auto {minutes} command will automatically send the exchange rate. The /stop_auto command will override /start_auto.`)
	})

	bot.Handle("/rates", func(c telebot.Context) error {
		var message string
		tag := c.Args()

		if len(tag) == 0 {
			currencies, err := service.GetCurrencies(context.Background())
			if err != nil {
				fmt.Printf("error in bot handle /rates: %v", err)
				return c.Send("Something wrong. Please try again later.")
			}

			for _, currency := range currencies {
				message += fmt.Sprintf("%s = %.2f ", currency.CurrencyName, currency.CurrencyPrice)
			}

			return c.Send(message)
		}

		if len(tag) == 1 {
			currency, err := service.GetCurrency(context.Background(), tag[0])
			if err != nil {
				return c.Send("Something wrong. Please try again later.")
			}

			message = fmt.Sprintf("%s = %.2f", currency.CurrencyName, currency.CurrencyPrice)

			return c.Send(message)
		}

		return c.Send("wrong arguments count")
	})

	bot.Handle("/start_auto", func(c telebot.Context) error {
		tag := c.Args()
		if len(tag) == 1 {
			n, err := strconv.Atoi(tag[0])
			if err != nil {
				return c.Send("Invalid parametr type. Only numbers.")
			}

			timeout = time.Duration(n) * time.Minute
			for {
				select {
				case <-autoChan:
					return nil

				case <-time.After(timeout):
					message, err := getCurrencies(service)
					if err != nil {
						fmt.Printf("error in getCurrencies: %v", err)
					}
					fmt.Println(message)
					c.Send(message)
				}
			}
		} else {
			return c.Send("Invalid parametrs count.")
		}
	})

	bot.Handle("/stop_auto", func(c telebot.Context) error {
		go func(chan struct{}) {
			autoChan <- struct{}{}
		}(autoChan)

		return c.Send("Autosender deactivated.")
	})

	bot.Start()
}
