package config

type Config struct {
	DataBase DataBase `yaml:"DataBase"`
	Host     Host     `yaml:"Host"`
	APIKey   string   `yaml:"APIKey"`
}

type DataBase struct {
	DBHost     string `yaml:"DBHost"`
	DBPort     string `yaml:"DBPort"`
	DBName     string `yaml:"DBName"`
	DBUser     string `yaml:"DBUser"`
	DBPassword string `yaml:"DBPassword"`
}

type Host struct {
	HostPort string `yaml:"HostPort"`
}
