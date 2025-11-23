package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"

	"github.com/yourusername/uas/config"
	"github.com/yourusername/uas/internal/middleware"
	pg "github.com/yourusername/uas/internal/repository/pg"
	mgo "github.com/yourusername/uas/internal/repository/mongo"
	"github.com/yourusername/uas/internal/route"
	"github.com/yourusername/uas/internal/service"
)

func main() {
	// load .env if present (no fatal if absent)
	_ = godotenv.Load()

	// load config from env
	cfg := config.LoadFromEnv()

	// ensure upload dir exists
	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		log.Fatalf("cannot create upload dir: %v", err)
	}

	// initialize Postgres pool
	pgPool, err := pg.NewPGPool(cfg.PGDsn)
	if err != nil {
		log.Fatalf("postgres connect error: %v", err)
	}
	// pgPool has Close(ctx) method, but pgxpool requires context for Close in some versions.
	defer func() {
		_ = pgPool.Close()
	}()

	// initialize Mongo client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoCli, err := mgo.NewMongoClient(ctx, cfg.MongoURI)
	if err != nil {
		log.Fatalf("mongo connect error: %v", err)
	}
	defer func() {
		_ = mongoCli.Disconnect(context.Background())
	}()

	// repositories
	userRepo := pg.NewUserRepoPG(pgPool)
	achRefRepo := pg.NewAchievementRefRepoPG(pgPool)
	achMongoRepo := mgo.NewAchievementRepoMongo(mongoCli.Database(cfg.MongoDB))

	// services (service methods will be used as Fiber handlers)
	authSvc := service.NewAuthService(userRepo, cfg)
	achSvc := service.NewAchievementService(achMongoRepo, achRefRepo, cfg.UploadDir)

	// fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	// global middleware: recover/panic handler
	app.Use(middleware.Recover())

	// Swagger (served at /swagger/index.html)
	app.Get("/swagger/*", swagger.HandlerDefault)

	// auth public routes and secured routes registration
	route.RegisterRoutes(app, authSvc, achSvc, cfg)

	// start server gracefully
	serverAddr := ":" + cfg.Port
	go func() {
		log.Printf("starting server on %s (env=%s)\n", serverAddr, os.Getenv("APP_ENV"))
		if err := app.Listen(serverAddr); err != nil {
			// if server closed normally, Listen returns an error; log if it's unexpected
			log.Fatalf("app.Listen error: %v", err)
		}
	}()

	// wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutdown signal received, shutting down...")

	// shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := app.Shutdown(); err != nil {
		log.Printf("fiber shutdown error: %v", err)
	}

	// close mongo & pg (some already deferred)
	if err := mongoCli.Disconnect(shutdownCtx); err != nil {
		log.Printf("mongo disconnect err: %v", err)
	}

	// pgxpool Close is synchronous; we already deferred Close above
	_ = pgPool.Close()

	log.Println("server stopped")
}
