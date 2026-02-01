package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/redis/go-redis/v9"

	"github.com/MagnunAVF/random-pass/internal/infra/repository"
	"github.com/MagnunAVF/random-pass/internal/service"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	repo := repository.NewRedisRepo(rdb)
	passwordSvc := service.NewPasswordService(repo)

	app := fiber.New(fiber.Config{
		AppName: "Random Pass v1",
	})

	app.Use(session.New())

	app.Get("/", func(c fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendFile("./index.html")
	})

	api := app.Group("/api")

	api.Post("/login", func(c fiber.Ctx) error {
		sess := session.FromContext(c)
		sess.Set("user_id", "dev")
		return c.JSON(fiber.Map{"status": "authenticated"})
	})

	api.Post("/logout", func(c fiber.Ctx) error {
		sess := session.FromContext(c)
		if err := sess.Destroy(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Logout failed")
		}
		return c.JSON(fiber.Map{"status": "logged_out"})
	})

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

		history, err := repo.GetLastFive(userID.(string))
		if err != nil {
			return c.JSON(fiber.Map{"history": []string{}})
		}

		return c.JSON(fiber.Map{"history": history})
	})

	log.Fatal(app.Listen(":3000"))
}
