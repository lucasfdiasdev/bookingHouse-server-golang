package main

import (
	"bookingHouses-server/routes"
	"bookingHouses-server/storage"

	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
)

func main() {
	godotenv.Load()
	storage.InitializeDB()

	app := iris.Default()

	location := app.Party("/api/location")
	{
		location.Get("/autocomplete", routes.Autocomplete)
		location.Get("/search", routes.Search)
	}
	user := app.Party("/api/user")
	{
		user.Post("/register", routes.Register)
	}

	app.Listen(":4000")
}
