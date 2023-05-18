package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"item-service/internal/app"
	"item-service/internal/repository"
	"item-service/internal/service"
	pb "item-service/pkg/gen/go/api/service/v1"
	"item-service/pkg/middleware"
	"log"
	"net/http"
)

func main() {
	// Prepare config file
	viper.AddConfigPath("../config")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}

	// Connect to a broker
	nc, err := nats.Connect(fmt.Sprintf("nats://%s:%d", viper.GetString("nats-host"), viper.GetInt("nats-port")))
	if err != nil {
		log.Fatalln(err)
	}

	// Connect to cache
	cache := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", viper.GetString("cache-host"), viper.GetInt("cache-port")),
		Password: viper.GetString("cache-password"),
		DB:       0, // use default DB
	})
	err = cache.Ping(context.Background()).Err()
	if err != nil {
		log.Fatalf("cannot ping db: %v", err)
	}

	// Connect to db
	db, err := repository.NewDB()
	if err != nil {
		log.Fatalf("cannot ping db: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("cannot ping db: %v", err)
	}

	// Register services
	dao := repository.NewDAO(cache)
	itemService := app.NewItemServiceServer(
		service.NewItemService(dao, nc),
	)

	// Configure cors policy
	c := cors.New(cors.Options{
		AllowedOrigins: viper.GetStringSlice("cors-origins"),
		AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE"},
	})

	// Open http gateway
	mux := runtime.NewServeMux(
		runtime.WithForwardResponseOption(middleware.HttpResponseModifier),
		runtime.WithErrorHandler(middleware.HandleRoutingError),
	)

	err = pb.RegisterItemServiceHandlerServer(context.Background(), mux, itemService)
	if err != nil {
		log.Fatalln("cannot register this service")
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", viper.GetString("http-hostname"), viper.GetInt("dev-http-handlers-port")),
		Handler: c.Handler(mux),
	}
	log.Println("Starting http gateway at ", viper.GetInt("dev-http-handlers-port"))

	log.Printf("Service %s started successfully\n", viper.GetString("service-name"))
	log.Fatalln(srv.ListenAndServe())
}
