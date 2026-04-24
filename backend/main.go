package main

import (
	"backend/config"
	"backend/handlers"
	"context"
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	defaultEnvironment = "production"
	corsMaxAge         = 12 * 3600
)

var cfg *config.Config

func init() {
	cfg = loadConfiguration()
}

func main() {
	ctx := createContext()
	defer ctx.Value("cancel").(context.CancelFunc)()

	startPprofServer()

	router := setupRouter(ctx)
	startServer(router)
}

func loadConfiguration() *config.Config {
	env := flag.String("env", defaultEnvironment, "Set the environment (development|production)")
	flag.Parse()

	configuration, err := config.LoadConfigForEnvironment(*env)
	if err != nil {
		log.Fatalf("Error loading configuration for environment %s: %v", *env, err)
	}

	log.Printf("Loaded configuration for environment: %s", *env)
	config.PrintConfig(configuration, *env)

	return configuration
}

func createContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	return context.WithValue(ctx, "cancel", cancel)
}

func startPprofServer() {
	runtime.MemProfileRate = 1 // Enable memory profiling

	pprofAddr := ":6060"
	go func() {
		log.Printf("Starting pprof server at %s", pprofAddr)
		log.Printf("Access profiling at:")
		log.Printf("  - Heap profile: http://localhost:6060/debug/pprof/heap")
		log.Printf("  - Goroutine profile: http://localhost:6060/debug/pprof/goroutine")
		log.Printf("  - CPU profile: http://localhost:6060/debug/pprof/profile?seconds=30")
		log.Printf("  - All profiles: http://localhost:6060/debug/pprof/")

		if err := http.ListenAndServe(pprofAddr, nil); err != nil {
			log.Printf("pprof server error: %v", err)
		}
	}()
}

func setupRouter(ctx context.Context) *gin.Engine {
	router := gin.Default()

	configureCORS(router)
	registerRoutes(router, ctx)

	return router
}

func configureCORS(router *gin.Engine) {
	corsConfig := cors.Config{
		AllowOrigins:     cfg.GetCORSConfig().AllowedOrigins,
		AllowMethods:     cfg.GetCORSConfig().AllowedMethods,
		AllowHeaders:     cfg.GetCORSConfig().AllowedHeaders,
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           corsMaxAge,
	}

	log.Printf("CORS configured with origins: %v", cfg.GetCORSConfig().AllowedOrigins)
	router.Use(cors.New(corsConfig))
}

func registerRoutes(router *gin.Engine, ctx context.Context) {
	handle := handlers.NewHandle(ctx, cfg.GetStaticFolder(), cfg.GetBackupFolder())
	handle.RegisterRoutes(router)
}

func startServer(router *gin.Engine) {
	serverAddress := cfg.GetServerAddress()
	log.Printf("Server started at %s", serverAddress)

	if err := router.Run(serverAddress); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
