package config

type Config struct {
	DBHost     string  `env:"DBHOST" envDefault:"localhost"`
	DBPort     string  `env:"DBPORT" envDefault:"5432"`
	Username string  `env:"DBUSER" `
	Password string  `env:"DBPASSWORD"`
	DBName   string  `env:"DBNAME"`
	SSLMode  string  `env:"SSLMODE" envDefault:"disable"`
	DOMAIN   string  `env:"DOMAIN"`
	APPPORT string  `env:"APPPORT" envDefault:"3000"`
	APIQUOTA int     `env:"APIQUOTA" envDefault:"10"`
}