package main

import (
	"log"
	"movie_reserve/database"
	"movie_reserve/router"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {

	database.InitDB()
	defer database.DB.Close()
	route := gin.Default()
	router.RegisterRoutes(route)
	err := route.Run(":8080")
	if err != nil {
		log.Fatalf("Error in porting %v:", err)
		return
	}

}
