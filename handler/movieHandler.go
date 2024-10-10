package handler

import (
	"movie_reserve/database"
	"movie_reserve/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetMovies(c *gin.Context) {
	var movies []model.Movie

	rows, err := database.DB.Query("SELECT id, title, description, director, duration, available_seats FROM movies")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying the data"})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var movie model.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Director, &movie.Duration, &movie.AvailableSeats); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning the rows"})
			return
		}
		movies = append(movies, movie)
	}
	c.JSON(http.StatusOK, movies)

}

func GetMovieByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie id"})
		return
	}
	var movie model.Movie
	err = database.DB.QueryRow("SELECT id, title, description, director, duration, available_seats FROM movies WHERE id = $1", id).
		Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Director, &movie.Duration, &movie.AvailableSeats)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Movie not found"})
		return
	}

	c.JSON(http.StatusOK, movie)
}

func SeatAvailability(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting movie id"})
		return
	}

	var availableSeats int

	err = database.DB.QueryRow("SELECT available_seats FROM movies Where id = $1", id).Scan(&availableSeats)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting available seats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"The available seats are": availableSeats})
}
