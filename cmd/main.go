package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"api/internal/handlers"
	"api/internal/models"
	utils "api/pkg"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.MaxMultipartMemory = 100 << 20
	var err error
	dsn := "host=localhost user=postgres dbname=moviesdb port=5432 sslmode=disable"
	utils.Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка создания базы данных: \n%v", err)
		return
	}

	if err := utils.Db.AutoMigrate(&models.Movie{}, &models.Member{}, &models.Genre{}, &models.MovieMember{}); err != nil {
		log.Fatal("Ошибка миграции: ", err)
	}

	r.Static("/media", "./media")

	r.GET("/movies", handlers.GetMovies)
	r.POST("/movies", handlers.AddMovie)
	r.GET("/movies/:id", handlers.GetMovie)
	r.DELETE("/movies/:id", handlers.DeleteMovie)
	r.PATCH("/movies/:id", handlers.UpdateMovie)

	r.GET("/genres", handlers.GetGenres)
	r.POST("/genres", handlers.AddGenre)
	r.DELETE("/genres/:id", handlers.DeleteGenre)
	r.PUT("/genres/:id", handlers.UpdateGenre)

	r.GET("/members", handlers.GetMembers)
	r.POST("/members", handlers.AddMember)
	r.GET("/members/:id", handlers.GetMember)
	r.DELETE("/members/:id", handlers.DeleteMember)
	r.PATCH("/members/:id", handlers.UpdateMember)

	r.Run(":8080")
}
