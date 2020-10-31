package main

import (
	"github.com/go-dawn/dawn"
	"github.com/go-dawn/dawn/config"
	"github.com/go-dawn/dawn/db/redis"
	"github.com/go-dawn/dawn/db/sql"
	"github.com/go-dawn/dawn/log"
)

func main() {
	config.Load("./_examples/application")

	sloop := dawn.New().
		AddModulers(
			log.New(nil),
			sql.New(),
			redis.New(),
			// add custom module
		)

	defer sloop.Cleanup()

	sloop.Setup().Watch()
}
