package archiver

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/rijdendetreinen/gotrain/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var redisDb *redis.Client

// Connect initializes the Redis client
func Connect() error {
	redisAddress := viper.GetString("archive.address")
	redisPassword := viper.GetString("archive.password")
	redisDbNumber := viper.GetInt("archive.db")

	log.WithField("address", redisAddress).
		WithField("password", redisPassword).
		WithField("db", redisDbNumber).
		Info("Connecting to Redis server")

	redisDb = redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: redisPassword,
		DB:       redisDbNumber,
	})

	result := redisDb.Ping()

	return result.Err()
}

// ProcessService adds a service object to the queue
func ProcessService(service models.Service) {
	serviceJSON, _ := json.Marshal(&service)

	fmt.Println(string(serviceJSON))

	if serviceJSON != nil {
		redisDb.LPush("services-queue", string(serviceJSON))
	}
}
