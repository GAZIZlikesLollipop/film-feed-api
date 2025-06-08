package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
	var movie Movie
	if err := db.First(&movie, id).Error; err != nil {
		log.Println("Фильм не найден: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Фильм не найден"})
		return
	}
	if err := db.Model(&movie).Association("Members").Clear(); err != nil {
		log.Printf("Ошибка очистки связей участников: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка очистки связей участников"})
		return
	}
	if err := db.Model(&movie).Association("Categories").Clear(); err != nil {
		log.Printf("Ошибка очистки связей категорий: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка очистки связей категорий"})
		return
	}
	if err := db.Delete(&movie).Error; err != nil {
		log.Println("Ошибка удаления фильма: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления фильма"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Фильм успешно удален"})
}

func updateMovie(c *gin.Context) {
	id := c.Param("id")
	var movie Movie
	if err := db.Preload("Categories").Preload("Members").First(&movie, id).Error; err != nil {
		log.Println("Фильм не найден", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Фильм не найден"})
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
