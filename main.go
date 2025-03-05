package main

import (
	"log"
	"database/sql"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"

	"github.com/Liptor/song_library/internal/config"
	"github.com/Liptor/song_library/handlers"
	"github.com/gofiber/fiber/v2"
)


func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file fount")
	}
}

func main() {
	app := fiber.New()
	dbconfig := config.New()
	dbinfo := "host=" + dbconfig.DB.Host + " port=" + dbconfig.DB.Port + " user=" + dbconfig.DB.User + " password=" + dbconfig.DB.Password + " dbname=" + dbconfig.DB.Name + " sslmode=disable"

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}

	defer db.Close()

	h := handlers.NewHandler(db)

	app.Post("/song", h.CreateSongHandler)
	app.Put("/song/:id", h.UpdataSongHandler)
	app.Delete("/song/:id", h.DeleteSongHandler)
	app.Get("/song", h.GetSongHandler)

	app.Listen(":3030")
}