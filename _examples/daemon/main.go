package main

import (
	"flag"

	"github.com/go-dawn/dawn"
	"github.com/go-dawn/dawn/config"
	"github.com/go-dawn/dawn/fiberx"
	"github.com/go-dawn/dawn/log"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// 🌶️ Notice that go run won't work in daemon mode
	// 🌶️ Please at dawn root dir and run go build -o play ./_examples/daemon
	// 🌶️ And run ./play
	config.Load("./_examples/daemon")
	config.LoadEnv()

	log.InitFlags(nil)
	flag.Parse()
	defer log.Flush()

	sloop := dawn.Default()

	router := sloop.Router()
	router.Get("/", func(c *fiber.Ctx) error {
		return fiberx.Message(c, "I'm running in daemon 🍀")
	})

	log.Infoln(0, sloop.Run(":3000"))
}
