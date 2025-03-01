package handlers

import (
	"database/sql"
	"log"
	"os"

	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	DB *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{DB: db}
}


func (h *Handler) CreateSongHandler(c *fiber.Ctx) error  {
	log.Println("Raw request body:", string(c.Body()))
	path, _ := os.LookupEnv("SONG_URL")

	var songdata struct {
		Group 	string `json:"group"`
		Song 	string	`json:"song"`
		ReleaseDate string `json:"releaseDate"`
		Text 		string `json:"text"`
		Link		string `json:"link"`
	}

	if err := c.BodyParser(&songdata); err != nil {
		return err
	}
	
	params := url.Values{}
	params.Add("group", songdata.Group)
	params.Add("song", songdata.Song)

	log.Println(path, "Path")
	fullURL := path + "?" + params.Encode()

	resp, err := http.Get(fullURL)
	if err != nil {
		log.Println("Error fetching song data:", err)
		return err
	}
	defer resp.Body.Close()

	if err != nil {
		log.Println("Error reading answer:", err)
		return err
	}	

	
	dbquery := `INSERT INTO songs ("group", song, "releaseDate", text, link) VALUES ($1, $2, $3, $4, $5)`

	_, err = h.DB.Exec(dbquery, 
		songdata.Group, 
		songdata.Song, 
		songdata.ReleaseDate, 
		songdata.Text, 
		songdata.Link,
	)
	
	if err != nil {
		log.Println("Error inserting quote into database:", err)
		return err
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "Song added successfully",
		"data": songdata,
	})
}

