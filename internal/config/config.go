package config

type Config struct {
	host     string  `env:"HOST" envDefault:"localhost"`
	port     string  `env:"PORT" envDefault:"3000"`
	username string  `env:"USERNAME" `
	password string  `env:"PASSWORD"`
	dbname   string  `env:"DBNAME"`
	sslmode  string  `env:"SSLMODE" envDefault:"false"`
}