package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DataBase DataBase `yaml:"dataBase"`
	Host     Host     `yaml:"host"`
	APIKey   string   `yaml:"apiKey"`
	TimeOut  int      `yaml:"timeOut"`
}

type DataBase struct {
	DBHost     string `yaml:"dbHost"`
	DBPort     string `yaml:"dbPort"`
	DBName     string `yaml:"dbName"`
	DBUser     string `yaml:"dbUser"`
	DBPassword string `yaml:"dbPassword"`
}

type Host struct {
	HostPort string `yaml:"hostPort"`
}

func NewConfig() (*Config, error) {
	var config Config

	configFile, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal config file: %w", err)
	}

	return &config, nil
}

func GetDSN(config *Config) string {
	dsn := "postgres://" + config.DataBase.DBUser + ":" + config.DataBase.DBPassword + "@" +
		config.DataBase.DBHost + ":" + config.DataBase.DBPort + "/" +
		config.DataBase.DBName + "?sslmode=disable"

	return dsn
}
