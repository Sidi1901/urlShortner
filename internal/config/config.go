package config

type Config struct {
	DBHost        string `env:"DBHOST" envDefault:"localhost"`
	DBPort        string `env:"DBPORT" envDefault:"5432"`
	Username      string `env:"DBUSER" `
	Password      string `env:"DBPASSWORD"`
	DBName        string `env:"DBNAME"`
	SSLMode       string `env:"SSLMODE" envDefault:"disable"`
	Domain        string `env:"DOMAIN"`
	AppPort       string `env:"APPPORT" envDefault:"3000"`
	APIQuota      int    `env:"APIQUOTA" envDefault:"10"`
	RedisDBNo     int    `env:"REDISDBNO"`
	RedisPassword string `env:"REDISPASSWORD"`
	RedisAddr     string `env:"REDISADDR" envDefault:"localhost:6379"`
	JwtSecret     string `env:"JWTSECRET"`
}
