package main

import (
	"log"
	"net/http"
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

func getMovies(c *gin.Context) {
	var movies []Movie
	if err := db.Preload("Categories").Preload("Members").Find(&movies).Error; err != nil {
		log.Println("Ошибка получения фильмов: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения фильмов"})
		return
	}
	c.JSON(http.StatusOK, movies)
}

func addMovie(c *gin.Context) {
	var movie Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		log.Println("Введены не верные данные", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Введены не верные данные"})
		return
	}
	if err := db.Create(&movie).Error; err != nil {
		log.Println("Ошибка добавленя фильма", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка добавелния фильма"})
		return
	}
	c.JSON(http.StatusCreated, movie)
}

func getMovie(c *gin.Context) {
	id := c.Param("id")
	var movie Movie
	if err := db.Preload("Members").Preload("Categories").Find(&movie, id); err != nil {
		log.Println("Фильм не найден", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Фильм не найден"})
		return
	}
	c.JSON(http.StatusOK, movie)
}

func deleteMovie(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&Movie{}, id).Error; err != nil {
		log.Println("Фильм не найден", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Фильм не найден"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Фильм успешно удален"})
}

func getIDs(data interface{}) []int64 {
	var ids []int64
	if items, ok := data.([]interface{}); ok {
		for _, item := range items {
			if m, ok := item.(map[string]interface{}); ok {
				if id, ok := m["id"].(float64); ok {
					ids = append(ids, int64(id))
				}
			}
		}
	}
	return ids
}

func updateMovie(c *gin.Context) {
	id := c.Param("id")
	var movie Movie
	if err := db.First(&movie, id).Error; err != nil {
		log.Println("Фильм не найден", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Фильм не найден"})
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		log.Println("Введены не верные данные: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Введены неверные данные"})
		return
	}
	if err := db.Model(&movie).Updates(updates).Error; err != nil {
		log.Println("Ошибка обнвления фильма: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления фильма"})
		return
	}
	if categories, ok := updates["categories"]; ok {
		var newCategories []Category
		if err := db.Where("id IN ?", getIDs(categories)).Find(&newCategories); err != nil {
			log.Println("Неверные данные: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверные данные"})
			return
		}
		if err := db.Model(&movie).Association("Categories").Replace(newCategories); err != nil {
			log.Println("Ошибка обновления категорий: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления категорий"})
			return
		}
	}
	if members, ok := updates["members"]; ok {
		var newMembers []Member
		if err := db.Where("id IN ?", getIDs(members)).Find(&newMembers); err != nil {
			log.Println("Неверные данные: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверные данные"})
			return
		}
		if err := db.Model(&movie).Association("Members").Replace(newMembers); err != nil {
			log.Println("Ошибка обновления участников: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления участников"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Фильм успешно обновлен", "movie": movie})
}

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

	r.GET("/movies", getMovies)
	r.POST("/movies", addMovie)
	r.GET("/movies/:id", getMovie)
	r.DELETE("/movies/:id", deleteMovie)
	r.PATCH("/movies/:id", updateMovie)

	r.Run(":8080")
}
