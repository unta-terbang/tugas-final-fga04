package main

import (
	"main.go/database"
	"main.go/router"
)

func main() {
	database.StartDB()

	router.Routers().Run(":8080")
}