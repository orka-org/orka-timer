package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/orka-org/orka-timer/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal("Could not connect to mongo: ", err)
	}
	defer client.Disconnect(context.Background())

	db := NewMongo(client, "orka", "timers")
	l, _ := zap.NewDevelopment()

	handler := NewHandler(db, l)
	_ = handler

	app := fiber.New()

	app.Post("/timers", handler.CreateTimer)
	app.Get("/timers", handler.ListTimers)
	app.Get("/timers/:id", handler.GetTimer)
	app.Put("/timers/:id", handler.UpdateTimer)
	app.Delete("/timers/:id", handler.DeleteTimer)

	// Start server
	port := getEnv("API_PORT", "4000")
	log.Fatal(app.Listen(":" + port))
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
