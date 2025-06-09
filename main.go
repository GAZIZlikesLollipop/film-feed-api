package main

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type Category struct {
	Id   int64  `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

type Movie struct {
	Id          int64      `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name"`
	PosterURL   string     `json:"posterURL"`
	MovieURL    string     `json:"movieURL"`
	Duration    int        `json:"duration"`
	Age         int        `json:"age"`
	Categories  []Category `json:"categories" gorm:"many2many:movie_categories;"`
	Rating      float32    `json:"rating"`
	Reviews     int64      `json:"reviews"`
	Description string     `json:"description"`
	Country     string     `json:"country"`
	Year        int        `json:"year"`
	Members     []Member   `json:"members" gorm:"many2many:movie_members;"`
	//Mb
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Member struct {
	Id            int64   `json:"id" gorm:"primaryKey"`
	Name          string  `json:"name"`
	FeaturedFilms []Movie `json:"featuredFilms" gorm:"many2many:movie_members;"`
	Photo         string  `json:"photo"`
	Role          string  `json:"role"`
	Date          string  `json:"date"`
	//Mb
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

var db *gorm.DB

func main() {
	r := gin.Default()
	var err error
	dsn := "host=localhost user=postgres dbname=moviesdb port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка создания базы данных: \n%v", err)
		return
	}

	if err := db.AutoMigrate(&Movie{}, &Member{}, &Category{}); err != nil {
		log.Fatal("Ошибка миграции: ", err)
	}

	r.Static("/media", "./static")

	r.GET("/movies", getMovies)
	r.POST("/movies", addMovie)
	r.GET("/movies/:id", getMovie)
	r.DELETE("/movies/:id", deleteMovie)
	r.PATCH("/movies/:id", updateMovie)

	r.GET("/categories", getCategories)
	r.POST("/categories", addCategory)
	r.DELETE("/categories/:id", deleteCategory)
	r.PUT("/categories/:id", updateCategory)

	r.Run(":8080")
}
