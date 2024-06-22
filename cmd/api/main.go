package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"user_news_api/handler"
	"user_news_api/notifier"
	"user_news_api/ratelimiter"
	"user_news_api/services"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

func main() {
	notifierOptions := getNotifierOptions()
	userNotifier := notifier.NewClient(notifierOptions)

	redisOptions := getRedisOptions()
	redisClient := redis.NewClient(redisOptions)
	limiter := ratelimiter.NewLimiterPool(redisClient, ratelimiter.DefaultConfigs)

	serv := services.NewUserNotifier(limiter, userNotifier)

	router := chi.NewRouter()

	handler.SetUserController(router, serv)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Fatal(server.ListenAndServe())
}

func getNotifierOptions() notifier.Options {
	host := os.Getenv("NOTIFIER_HOST")
	if host == "" {
		panic("notifier host is empty")
	}

	portStr := os.Getenv("NOTIFIER_PORT")
	if portStr == "" {
		panic("notifier port is empty")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic("notifier port is not a number")
	}

	sender := os.Getenv("NOTIFIER_SENDER")
	if sender == "" {
		panic("notifier sender is empty")
	}

	password := os.Getenv("NOTIFIER_PASSWORD")
	if password == "" {
		panic("notifier password is empty")
	}

	return notifier.Options{
		Host:     host,
		Port:     port,
		Username: sender,
		Password: password,
	}
}

func getRedisOptions() *redis.Options {
	addr := os.Getenv("REDIS_ADDRESS")
	if addr == "" {
		panic("redis address is empty")
	}

	// password is not validated due to local environment use case.
	password := os.Getenv("REDIS_PASSWORD")

	return &redis.Options{
		Addr:     addr,
		Password: password,
	}
}
