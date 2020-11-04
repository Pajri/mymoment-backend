package config

type Configuration struct {
	DB                DBConfig
	SMTP              SMTP
	Host              string
	FEHost            string
	EmailVerification EmailVerificationConfig
	ResetPassword     ResetPasswordConfig
	Redis             RedisConfig
}

type DBConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DbName   string
}

type SMTP struct {
	From     string
	Host     string
	Port     int
	Username string
	Password string
}

type EmailVerificationConfig struct {
	Subject string
}

type ResetPasswordConfig struct {
	Subject string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
}
