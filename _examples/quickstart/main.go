package main

import (
	"log"

	"github.com/go-dawn/dawn"
	"github.com/go-dawn/dawn/fiberx"
	"github.com/gofiber/fiber/v2"
)

func main() {
	sloop := dawn.Default()

	router := sloop.Router()
	// GET /  =>  Welcome to dawn 👋
	router.Get("/", func(c *fiber.Ctx) error {
		return fiberx.Message(c, "Welcome to dawn 👋")
	})

	log.Println(sloop.Run(":3000"))
}
