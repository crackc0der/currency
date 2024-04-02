package main

import (
	"exchange_course/config"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func Run() {
	configFile, err := os.ReadFile("../../config/config.yaml")

	if err != nil {
		log.Fatal("Could not read config file: ", err)
	}
	fmt.Println(configFile)
	config := config.Config{}
	err = yaml.Unmarshal(configFile, &config)

	if err != nil {
		log.Fatal("Could not unmarshal config file: ", err)
	}
	fmt.Println(config.DataBase.DbName)
}
