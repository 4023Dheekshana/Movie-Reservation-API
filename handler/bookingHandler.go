package handler

import (
	"fmt"
	"log"
	"movie_reserve/database"
	"movie_reserve/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateBooking(c *gin.Context) {
	var booking model.Booking

	if c.Request.Header.Get("Content-Type") != "application/json" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Content-Type"})
		fmt.Println("Invalid Content-Type")
		return
	}

	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error binding json"})
		return
	}

	var availableSeats int

	err := database.DB.QueryRow("SELECT available_seats FROM movies WHERE id = $1", booking.MovieID).Scan(&availableSeats)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Movie not found"})
		return
	}

	if booking.Seats > availableSeats {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Seats not available"})
		return
	}

	_, err = database.DB.Exec("INSERT INTO bookings (user_id, movie_id, seats, status) VALUES($1, $2, $3, $4)",
		booking.UserID, booking.MovieID, booking.Seats, booking.Status)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error executing the error"})
		log.Fatalf("error in executing: %v", err)
		return
	}

	_, err = database.DB.Exec("UPDATE movies SET available_seats = available_seats - $1 WHERE id = $2", booking.Seats, booking.MovieID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating the seat availability"})
		return
	}

	var bookingId int

	err = database.DB.QueryRow("SELECT id FROM bookings WHERE user_id = $1", booking.UserID).Scan(&bookingId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking created successfully"})
	c.JSON(http.StatusOK, gin.H{"Your booking id is": bookingId})

}

func GetBooking(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking id"})
	}
	var booking model.Booking

	err = database.DB.QueryRow("SELECT id, user_id, movie_id, seats, status FROM bookings WHERE id = $1", id).
		Scan(&booking.ID, &booking.UserID, &booking.MovieID, &booking.Seats, &booking.Status)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking id not found"})
	}
	c.JSON(http.StatusOK, booking)

}

func CancelBooking(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking id"})
		return
	}
	_, err = database.DB.Exec("DELETE FROM bookings WHERE id = $1", id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error deleting booking id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking cancelled successfully"})

}

func Payment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Booking id"})
		return
	}
	var booking model.Booking
	booking.Payment = "PAID"
	_, err = database.DB.Exec("UPDATE bookings SET payment = $1 WHERE id = $2", booking.Payment, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error updating payment"})
		log.Fatalf("error in updating %v:", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Payment paid successfully"})
}
