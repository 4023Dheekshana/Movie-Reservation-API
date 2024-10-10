package router

import (
	"movie_reserve/handler"
	"movie_reserve/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(route *gin.Engine) {

	auth := route.Group("/")
	auth.Use(middleware.Authenticate)

	route.POST("/signup", handler.UserSignUp)
	route.POST("/login", handler.Userlogin)

	route.GET("/getmovies", handler.GetMovies)
	route.GET("/getmovie/:id", handler.GetMovieByID)

	auth.GET("/availableseats/:id", handler.SeatAvailability)
	auth.POST("/bookmovie", handler.CreateBooking)
	auth.GET("/getbooking/:id", handler.GetBooking)
	auth.DELETE("/cancelticket/:id", handler.CancelBooking)
	auth.POST("/payment/:id", handler.Payment)

}
