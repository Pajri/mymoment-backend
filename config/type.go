package config

type Configuration struct {
	DB DBConfig
}

type DBConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DbName   string
}
