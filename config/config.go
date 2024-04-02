package config

type Config struct {
	DataBase DataBase `yaml:"DataBase"`
	Host     Host     `yaml:"Host"`
	ApiKey   string   `yamk:"ApiKey"`
}

type DataBase struct {
	DbHost     string `yaml:"DbHost"`
	DbPort     string `yaml:"DbPort"`
	DbName     string `yaml:"DbName"`
	DbUser     string `yaml:"DbUser"`
	DbPassword string `yaml:"DbPassword"`
}

type Host struct {
	HostPort string `yaml:"HostPort"`
}
