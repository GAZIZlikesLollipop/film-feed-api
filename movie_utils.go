package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	name := c.PostForm("name")
	duration, _ := strconv.Atoi(c.PostForm("duration"))
	age, _ := strconv.Atoi(c.PostForm("age"))
	rating, _ := strconv.ParseFloat(c.PostForm("rating"), 32)
	reviews, _ := strconv.ParseInt(c.PostForm("reviews"), 10, 64)
	description := c.PostForm("description")
	country := c.PostForm("country")
	year, _ := strconv.Atoi(c.PostForm("year"))

	categoryStr := c.PostForm("categories")
	var categories []Category
	if categoryStr != "" {
		ids := strings.Split(categoryStr, ",")
		for _, idStr := range ids {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Printf("Неверный ID категории: %s", idStr)
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Неверный ID категории: %s", idStr)})
				return
			}
			var cat Category
			if err := db.First(&cat, id).Error; err != nil {
				log.Printf("Категория %d не найдена: %v", id, err)
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Категория %d не найдена", id)})
				return
			}
			categories = append(categories, cat)
		}
	}

	memberStr := c.PostForm("members")
	var members []Member
	if memberStr != "" {
		ids := strings.Split(memberStr, ",")
		for _, idStr := range ids {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Printf("Неверный ID участника: %s", idStr)
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Неверный ID участника: %s", idStr)})
				return
			}
			var member Member
			if err := db.First(&member, id).Error; err != nil {
				log.Printf("Участник %d не найдена: %v", id, err)
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Участник %d не найдена", id)})
				return
			}
			members = append(members, member)
		}
	}

	posterFile, err := c.FormFile("posterURL")
	if err != nil {
		log.Println("Ошибка получения файла: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл не загружен"})
		return
	}
	movieFile, err := c.FormFile("movieURL")
	if err != nil {
		log.Println("Ошибка получения файла: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл не загружен"})
		return
	}

	movieDir := "static/movie-videos"
	if err := os.MkdirAll(movieDir, 0755); err != nil {
		log.Printf("Ошибка создания папки %s: %v", movieDir, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}

	movieFileName := fmt.Sprintf("%s-%s%s", name, uuid.New(), strings.ToLower(filepath.Ext(movieFile.Filename)))
	moviePath := filepath.Join(movieDir, movieFileName)

	if err := c.SaveUploadedFile(movieFile, moviePath); err != nil {
		log.Printf("Ошибка сохранения файла: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения файла"})
		return
	}

	posterDir := "static/movie-posters"
	if err := os.MkdirAll(posterDir, 0755); err != nil {
		log.Printf("Ошибка создания папки %s: %v", posterDir, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}

	posterFileName := fmt.Sprintf("%s-%s%s", name, uuid.New(), strings.ToLower(filepath.Ext(posterFile.Filename)))
	posterPath := filepath.Join(posterDir, posterFileName)

	if err := c.SaveUploadedFile(posterFile, posterPath); err != nil {
		log.Printf("Ошибка сохранения файла: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения файла"})
		return
	}

	movie := Movie{
		Name:        name,
		PosterURL:   fmt.Sprintf("http://192.168.1.9:8080/%s", posterPath),
		MovieURL:    fmt.Sprintf("http://192.168.1.9:8080/%s", moviePath),
		Duration:    duration,
		Age:         age,
		Rating:      float32(rating),
		Reviews:     reviews,
		Description: description,
		Country:     country,
		Year:        year,
		Categories:  categories,
		Members:     members,
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
