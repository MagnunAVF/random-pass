package main

import (
	"context"
	"embed"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/redis/go-redis/v9"

	"github.com/MagnunAVF/random-pass/internal/handler"
	infradb "github.com/MagnunAVF/random-pass/internal/infra/db"
	"github.com/MagnunAVF/random-pass/internal/infra/repository"
	"github.com/MagnunAVF/random-pass/internal/service"
)

//go:embed db/migrations/*.sql
var migrationsFS embed.FS

func main() {
	ctx := context.Background()

	pgPool, err := infradb.NewPool(ctx, getEnv("POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/randompass?sslmode=disable"))
	if err != nil {
		log.Fatalf("connect to postgres: %v", err)
	}
	defer pgPool.Close()

	if err := infradb.RunMigrations(ctx, pgPool, migrationsFS); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: getEnv("REDIS_ADDR", "localhost:6379"),
	})

	passwordRepo := repository.NewRedisRepo(rdb)
	userRepo := repository.NewPostgresUserRepo(pgPool)

	passwordSvc := service.NewPasswordService(passwordRepo)
	authSvc := service.NewAuthService(userRepo)

	authHandler := handler.NewAuthHandler(authSvc)

	app := fiber.New(fiber.Config{
		AppName: "Random Pass v1",
	})

	app.Use(session.New())

	app.Get("/", func(c fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendFile("./index.html")
	})

	api := app.Group("/api")

	api.Post("/signup", authHandler.Signup)
	api.Post("/login", authHandler.Login)
	api.Post("/logout", authHandler.Logout)

	api.Post("/generate", func(c fiber.Ctx) error {
		sess := session.FromContext(c)
		userID := sess.Get("user_id")

		if userID == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Please login first"})
		}

		pwd, err := passwordSvc.GenerateAndSave(userID.(string))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Generation failed"})
		}

		return c.JSON(fiber.Map{"password": pwd})
	})

	api.Get("/history", func(c fiber.Ctx) error {
		sess := session.FromContext(c)
		userID := sess.Get("user_id")

		if userID == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		history, err := passwordRepo.GetLastFive(userID.(string))
		if err != nil {
			return c.JSON(fiber.Map{"history": []string{}})
		}

		return c.JSON(fiber.Map{"history": history})
	})

	log.Fatal(app.Listen(":3000"))
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
