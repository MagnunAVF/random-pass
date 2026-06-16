package main

import (
	"context"
	"embed"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
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
	jwtSvc := service.NewJWTService(getEnv("JWT_SECRET", "change-me-in-production"))

	authHandler := handler.NewAuthHandler(authSvc, jwtSvc)

	app := fiber.New(fiber.Config{
		AppName: "Random Pass v1",
	})

	app.Get("/", func(c fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendFile("./index.html")
	})

	api := app.Group("/api")

	api.Post("/signup", authHandler.Signup)
	api.Post("/login", authHandler.Login)

	// JWT-protected routes
	protected := api.Group("", jwtMiddleware(jwtSvc))

	protected.Post("/generate", func(c fiber.Ctx) error {
		userID := c.Locals(handler.UserIDLocal).(string)

		pwd, err := passwordSvc.GenerateAndSave(userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Generation failed"})
		}

		return c.JSON(fiber.Map{"password": pwd})
	})

	protected.Get("/history", func(c fiber.Ctx) error {
		userID := c.Locals(handler.UserIDLocal).(string)

		history, err := passwordRepo.GetLastFive(userID)
		if err != nil {
			return c.JSON(fiber.Map{"history": []string{}})
		}

		return c.JSON(fiber.Map{"history": history})
	})

	log.Fatal(app.Listen(":3000"))
}

func jwtMiddleware(jwtSvc service.JWTServicer) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing or invalid authorization header"})
		}

		claims, err := jwtSvc.Validate(strings.TrimPrefix(authHeader, "Bearer "))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
		}

		c.Locals(handler.UserIDLocal, claims.UserID)
		return c.Next()
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
