package config

type Config struct {
	DataBase DataBase `yaml:"dataBase"`
	Host     Host     `yaml:"host"`
	APIKey   string   `yaml:"apiKey"`
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
