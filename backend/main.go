package main

import (
	"backend/config"
	"backend/handlers"
	"flag"
	"log"

	_ "backend/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	httpSwagger "github.com/swaggo/http-swagger"
)

var cfg *config.Config

func init() {
	// Use the Flag manage env variables
	env := flag.String("env", "production", "Set the environment (development|production)")
	flag.Parse()

	// load config
	var config_err error
	cfg, config_err = config.LoadConfigForEnvironment(*env)
	if config_err != nil {
		log.Fatalf("Error loading configuration for environment %s: %v", *env, config_err)
	}
	log.Printf("Loaded configuration for environment: %s", *env)
	config.PrintConfig(cfg, *env)
}

func main() {
	handle := handlers.NewHandle(cfg.GetStaticFolder())

	r := gin.Default()
	r.Use(cors.Default())

	// API routes
	handle.RegisterRoutes(r)
	// Static route
	r.Static("/static", cfg.GetStaticFolder())
	// Swagger UI route
	r.GET("/swagger/*any", gin.WrapH(httpSwagger.WrapHandler))

	log.Printf("Server started at %s", cfg.GetServerAddress())
	if err := r.Run(cfg.GetServerAddress()); err != nil {
		log.Fatal(err)
	}
}
