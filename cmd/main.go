package main

import (
"context"
"fmt"
"log"
"os"
"os/signal"
"time"

"github.com/gofiber/fiber/v2"
"github.com/gofiber/swagger"
"github.com/joho/godotenv"

"uas-backend/config"
"uas-backend/internal/middleware"
pg "uas-backend/internal/repository/pg"
mgo "uas-backend/internal/repository/mongo"
"uas-backend/internal/route"
"uas-backend/internal/service"

)

func main() {
// Load .env file if present
_ = godotenv.Load()

// Load configuration
cfg := config.LoadConfig()

// Ensure upload directory exists
if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
	log.Fatalf("cannot create upload dir: %v", err)
}

// Initialize Postgres pool
pgDsn := fmt.Sprintf(
  "postgres://%s:%s@%s:%s/%s?sslmode=disable",
  cfg.PGUser,
  cfg.PGPassword,
  cfg.PGHost,
  cfg.PGPort,
  cfg.PGDBName,
)
pgPool, err := pg.NewPGPool(pgDsn)
if err != nil {
	log.Fatalf("postgres connect error: %v", err)
}
defer pgPool.Close()

// Initialize Mongo client
mongoCli, err := mgo.NewMongoClient(cfg.MongoURI)
if err != nil {
	log.Fatalf("mongo connect error: %v", err)
}
defer func() {
	_ = mongoCli.Disconnect(context.Background())
}()

// Repositories
userRepo := pg.NewUserRepo(pgPool)
studentRepo := pg.NewStudentRepo(pgPool)
lecturerRepo := pg.NewLecturerRepo(pgPool)
achRefRepo := pg.NewAchievementRefRepo(pgPool)
achMongoRepo := mgo.NewAchievementMongoRepo(mongoCli.Database(cfg.MongoDB))

// Services
authSvc := service.NewAuthService(userRepo, cfg)
userSvc := service.NewUserService(userRepo)
studentSvc := service.NewStudentService(studentRepo)
lecturerSvc := service.NewLecturerService(lecturerRepo)
achSvc := service.NewAchievementService(achRefRepo, achMongoRepo)

// Fiber app
app := fiber.New(fiber.Config{
	ReadTimeout:  15 * time.Second,
	WriteTimeout: 15 * time.Second,
	IdleTimeout:  60 * time.Second,
})

// Global middleware
app.Use(middleware.Recover())

// Swagger docs
app.Get("/swagger/*", swagger.HandlerDefault)

// Register routes
route.RegisterRoutes(
	app,
	authSvc,
	userSvc,
	studentSvc,
	lecturerSvc,
	achSvc,
	cfg.JWTSecret, 
)


// Start server gracefully
serverAddr := ":" + cfg.Port
go func() {
	log.Printf("starting server on %s (env=%s)\n", serverAddr, os.Getenv("APP_ENV"))
	if err := app.Listen(serverAddr); err != nil {
		log.Fatalf("app.Listen error: %v", err)
	}
}()

// Wait for interrupt signal for graceful shutdown
quit := make(chan os.Signal, 1)
signal.Notify(quit, os.Interrupt)
<-quit
log.Println("shutdown signal received, shutting down...")

shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
defer shutdownCancel()

if err := app.Shutdown(); err != nil {
	log.Printf("fiber shutdown error: %v", err)
}

if err := mongoCli.Disconnect(shutdownCtx); err != nil {
	log.Printf("mongo disconnect err: %v", err)
}

log.Println("server stopped")
}

