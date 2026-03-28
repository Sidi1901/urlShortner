package database

import(
	"context"
    "fmt"
	"os"
	"github.com/Sidi1901/urlShortner/internal/config"
	"github.com/redis/go-redis/v9"
)


func ConnectMemoryDB(cfg *config.Config) *redis.Client {
	var ctx := context.TimeOut(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("Connecting to Redis database...")

	rdb := redis.NewClient(&redis.Options{
		Addr : cfg.RedisAddr,
		Password : cfg.RedisPassword,
		DB : cfg.RedisDBNo,
	})


	// 	PingContext is used to check if the database connection is alive within the specified context timeout.
	if err := rdb.Ping(ctx).Err(); err != nil {
		fmt.Println("Redis Connection error: ", err)
		panic(err)
	}

	fmt.Println("Connected to Redis database...")
	return rdb
}
