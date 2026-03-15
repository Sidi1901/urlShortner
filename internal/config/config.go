package config

type Config struct {
	Host     string  `env:"HOST" envDefault:"localhost"`
	Port     string  `env:"PORT" envDefault:"3000"`
	Username string  `env:"USERNAME" `
	Password string  `env:"PASSWORD"`
	DBName   string  `env:"DBNAME"`
	SSLMode  string  `env:"SSLMODE" envDefault:"false"`
}