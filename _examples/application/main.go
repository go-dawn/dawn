package main

import (
	"flag"

	"github.com/go-dawn/dawn"
	"github.com/go-dawn/dawn/config"
	"github.com/go-dawn/dawn/db/redis"
	"github.com/go-dawn/dawn/db/sql"
	"github.com/go-dawn/dawn/log"
)

func main() {
	config.Load("./_examples/application")
	config.LoadEnv()

	log.InitFlags(nil)
	flag.Parse()
	defer log.Flush()

	sloop := dawn.New().
		AddModulers(
			sql.New(),
			redis.New(),
			// add custom module
		)

	defer sloop.Cleanup()

	sloop.Setup().Watch()
}
