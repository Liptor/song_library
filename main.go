package main

import (
	"log"
	"strconv"

	"database/sql"

	_ "github.com/lib/pq"

	"github.com/Liptor/song_library/handlers"
	"github.com/gofiber/fiber/v2"
)

const (
	user = "acer"
	port = 5432
	password = "password"
	host = "localhost"
	db_name = "song_library"
)

func main() {
	app := fiber.New()
	portStr := strconv.Itoa(port)
	dbinfo := "host=" + host + " port=" + portStr + " user=" + user + " password=" + password + " dbname=" + db_name + " sslmode=disable"

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}

	defer db.Close()

	h := handlers.NewHandler(db)


	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS songs (
                        id 			SERIAL PRIMARY KEY,
						releaseDate DATE,
						text 		TEXT,
						link		TEXT
                        "group" 		TEXT NOT NULL,
						song  		TEXT
						
                    )`)

					

	if err != nil {
		log.Fatal(err)

	}

	app.Post("/create", h.CreateSongHandler)
	// app.Put("/update/:id", handlers.updataSongHandler)
	// app.Delete("/delete/:id", handlers.deleteSongHandler)
	// app.Get("/song", handlers.getSongHandler)
	// app.Get("/song/pagin")

	app.Listen(":3030")
}