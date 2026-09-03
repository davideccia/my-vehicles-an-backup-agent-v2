package main

import (
	"my_vehicles_an_backup_agent/bootstrap"
)

func main() {
	app := bootstrap.Boot()

	app.Start()
}
