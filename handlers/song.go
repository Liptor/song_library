package handlers

import (
	"fmt"
	"database/sql"
	"log"
	"os"
	"strings"
	"strconv"


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

type Song struct {
		Group 	string `json:"group"`
		Song 	string	`json:"song"`
		ReleaseDate string `json:"releaseDate"`
		Text 		string `json:"text"`
		Link		string `json:"link"`
}


func (h *Handler) CreateSongHandler(c *fiber.Ctx) error  {
	songdata := new(Song)
	log.Println("Raw request body:", string(c.Body()))
	path, _ := os.LookupEnv("SONG_URL")

	if err := c.BodyParser(songdata); err != nil {
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


func (h *Handler) DeleteSongHandler(c *fiber.Ctx) error {
	res := c.Params("id")
	if res == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Id is necessary")
	}

	dbquery := `DELETE FROM songs WHERE id=$1`

	_, err := h.DB.Exec(dbquery, res)

	if err != nil {
		return err
	}


	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "Song deleted successfully",
	})
}

func (h *Handler) UpdataSongHandler(c *fiber.Ctx) error {
	var songdata map[string]interface{}
	id := c.Params("id")

	if err := c.BodyParser(&songdata); err != nil {
		return err
	}

	setSong := []string{}
	values := []interface{}{}
	paramCount := 1

	for field, value := range songdata {
		if value != "" {
			escapedField := fmt.Sprintf(`"%s"`, field)
			setSong = append(setSong, fmt.Sprintf("%s = $%d", escapedField, paramCount))
			values = append(values, value)

			paramCount++
		}
	}


	setQuery := strings.Join(setSong, ", ")

	dbquery := fmt.Sprintf("UPDATE songs SET %s WHERE id = $%d", setQuery, paramCount)

	values = append(values, id)

	_, err := h.DB.Exec(dbquery, values...)


	if err != nil {
		log.Println("Unable to add to database: ", err)
		return err
	}

	return nil

}

func (h *Handler) GetSongHandler(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page"))
	
	if err != nil || page < 1 {
		page = 1 // Default page value
	}

	limit, err := strconv.Atoi(c.Query("limit"))

	if err != nil || limit < 1 {
		limit = 10 // Default page value
	}

	offset := (page - 1) * limit

	rows, err := h.DB.Query(`SELECT id, "group", song, releasedate, text, link FROM songs LIMIT $1 OFFSET $2`, limit, offset)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch items")
	}

	defer rows.Close()

	items := []map[string]interface{}{} 
 
	for rows.Next() {
		var id int
		var group, song, text, link string
		var releasedate  string

		if err := rows.Scan(&id, &group, &releasedate, &song, &text, &link); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to scan items")
		}

		items = append(items, map[string]interface{}{
			"id": id,
			"group": group,
			"song": song,
			"releasedate": releasedate,
			"text": text,
			"link": link,
		})
	}

	return c.JSON(fiber.Map{
		"page": page,
		"limit": limit,
		"items": items,
	}, "application/json")

}


