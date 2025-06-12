package handlers

import (
	"api/internal/models"
	utils "api/pkg"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func AddMovie(c *gin.Context) {
	name := c.PostForm("name")
	duration, _ := strconv.Atoi(c.PostForm("duration"))
	age, _ := strconv.Atoi(c.PostForm("age"))
	rating, _ := strconv.ParseFloat(c.PostForm("rating"), 32)
	reviews, _ := strconv.ParseInt(c.PostForm("reviews"), 10, 64)
	description := c.PostForm("description")
	country := c.PostForm("country")
	year, _ := strconv.Atoi(c.PostForm("year"))
	budget, _ := strconv.ParseInt(c.PostForm("budget"), 10, 64)
	boxOffice, _ := strconv.ParseInt(c.PostForm("boxOffice"), 10, 64)

	moviePath, err := utils.SaveFile(c, "movieURL", name, "media/movie-videos")
	if err != nil {
		log.Println("Ошибка получения пути файла: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintln("Ошибка получения пути файла: ", err)})
		return
	}

	trailerPath, err := utils.SaveFile(c, "trailerURL", fmt.Sprintf("%s-trailer", name), "media/movie-trailers")
	if err != nil {
		log.Println("Ошибка получения пути файла: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintln("Ошибка получения пути файла: ", err)})
		return
	}

	posterPath, err := utils.SaveFile(c, "posterURL", name, "media/movie-posters")
	if err != nil {
		log.Println("Ошибка получения пути файла: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintln("Ошибка получения пути файла: ", err)})
		return
	}

	genresStr := c.PostForm("genres")
	var genres []models.Genre
	if genresStr != "" {
		ids := strings.Split(genresStr, ",")
		for _, idStr := range ids {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				log.Printf("Неверный ID жанры: %s", idStr)
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Неверный ID жанры: %s", idStr)})
				return
			}
			var cat models.Genre
			if err := utils.Db.First(&cat, id).Error; err != nil {
				log.Printf("Жанр %d не найдена: %v", id, err)
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Жанр %d не найдена", id)})
				return
			}
			genres = append(genres, cat)
		}
	} else {
		genres = []models.Genre{}
	}

	movieMemberStr := c.PostForm("movieMembers")
	var movieMembers []models.MovieMember
	if movieMemberStr != "" {
		if err := json.Unmarshal([]byte(movieMemberStr), &movieMembers); err != nil {
			log.Printf("Неверный формат movieMembers: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Неверный формат movieMembers: %v", err)})
			return
		}
	} else {
		movieMembers = []models.MovieMember{}
	}

	movie := models.Movie{
		Name:         name,
		PosterURL:    fmt.Sprintf("http://192.168.1.9:8080/%s", posterPath),
		MovieURL:     fmt.Sprintf("http://192.168.1.9:8080/%s", moviePath),
		TrailerURL:   fmt.Sprintf("http://192.168.1.9:8080/%s", trailerPath),
		Duration:     duration,
		Age:          age,
		Rating:       float32(rating),
		Reviews:      reviews,
		Description:  description,
		Country:      country,
		Year:         year,
		Genres:       genres,
		Budget:       budget,
		BoxOffice:    boxOffice,
		MovieMembers: movieMembers,
	}

	if err := utils.Db.Preload("Genres").Preload("MovieMembers.Member").Create(&movie).Error; err != nil {
		log.Println("Ошибка добавленя фильма", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка добавелния фильма"})
		return
	}

	c.JSON(http.StatusCreated, movie)
}

func GetMovies(c *gin.Context) {
	var movies []models.Movie
	if err := utils.Db.Preload("Genres").Preload("MovieMembers.Member").Find(&movies).Error; err != nil {
		log.Println("Ошибка получения фильмов: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения фильмов"})
		return
	}
	c.JSON(http.StatusOK, movies)
}

func GetMovie(c *gin.Context) {
	id := c.Param("id")
	var movie models.Movie
	if err := utils.Db.Preload("MovieMembers.Member").Preload("Genres").First(&movie, id); err != nil {
		log.Println("Фильм не найден", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Фильм не найден"})
		return
	}
	c.JSON(http.StatusOK, movie)
}

func DeleteMovie(c *gin.Context) {
	id := c.Param("id")
	var movie models.Movie
	if err := utils.Db.First(&movie, id).Error; err != nil {
		log.Println("Фильм не найден: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Фильм не найден"})
		return
	}
	absolutePath := "/home/lollipop/dev/film-feed"
	if movie.MovieURL != "" {
		file, err := url.Parse(movie.MovieURL)
		if err != nil {
			log.Println("Ошибка прасиинга Url", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintln("Ошиибка парсиинга URL", err)})
			return
		}

		if err := os.Remove(fmt.Sprintf("%s%s", absolutePath, file.Path)); err != nil {
			log.Println("Ошибка удаления файла: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintln("Ошибка удаления файла: ", err)})
			return
		}
	}
	if movie.PosterURL != "" {
		file, err := url.Parse(movie.PosterURL)
		if err != nil {
			log.Println("Ошибка прасиинга Url", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintln("Ошиибка парсиинга URL", err)})
			return
		}

		if err := os.Remove(fmt.Sprintf("%s%s", absolutePath, file.Path)); err != nil {
			log.Println("Ошибка удаления файла: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintln("Ошибка удаления файла: ", err)})
			return
		}
	}
	if movie.TrailerURL != "" {
		file, err := url.Parse(movie.TrailerURL)
		if err != nil {
			log.Println("Ошибка прасиинга Url", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintln("Ошиибка парсиинга URL", err)})
			return
		}

		if err := os.Remove(fmt.Sprintf("%s%s", absolutePath, file.Path)); err != nil {
			log.Println("Ошибка удаления файла: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintln("Ошибка удаления файла: ", err)})
			return
		}
	}
	if err := utils.Db.Model(&movie).Association("MovieMembers").Clear(); err != nil {
		log.Printf("Ошибка очистки связей участников: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка очистки связей участников"})
		return
	}
	if err := utils.Db.Model(&movie).Association("Genres").Clear(); err != nil {
		log.Printf("Ошибка очистки связей категорий: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка очистки связей категорий"})
		return
	}
	if err := utils.Db.Delete(&movie).Error; err != nil {
		log.Println("Ошибка удаления фильма: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления фильма"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Фильм успешно удален"})
}

func UpdateMovie(c *gin.Context) {
	id := c.Param("id")
	var movie models.Movie
	if err := utils.Db.Preload("Genres").Preload("MovieMembers.Member").First(&movie, id).Error; err != nil {
		log.Println("Фильм не найден", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Фильм не найден"})
		return
	}

	name := c.PostForm("name")
	if name != "" {
		movie.Name = name
	}
	if durationStr := c.PostForm("duration"); durationStr != "" {
		duration, err := strconv.Atoi(durationStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат длительности"})
			return
		}
		movie.Duration = duration
	}
	if ageStr := c.PostForm("age"); ageStr != "" {
		age, err := strconv.Atoi(ageStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат возрастного ограничения"})
			return
		}
		movie.Age = age
	}
	if ratingStr := c.PostForm("rating"); ratingStr != "" {
		rating, err := strconv.ParseFloat(ratingStr, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат рейтинга"})
			return
		}
		movie.Rating = float32(rating)
	}
	if reviewsStr := c.PostForm("reviews"); reviewsStr != "" {
		reviews, err := strconv.ParseInt(reviewsStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат количества отзывов"})
			return
		}
		movie.Reviews = reviews
	}
	if description := c.PostForm("description"); description != "" {
		movie.Description = description
	}
	if country := c.PostForm("country"); country != "" {
		movie.Country = country
	}
	if yearStr := c.PostForm("year"); yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат года"})
			return
		}
		movie.Year = year
	}
	if budgetStr := c.PostForm("budget"); budgetStr != "" {
		budget, err := strconv.ParseInt(budgetStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат бюджета"})
			return
		}
		movie.Budget = budget
	}
	if boxOfficeStr := c.PostForm("boxOffice"); boxOfficeStr != "" {
		boxOffice, err := strconv.ParseInt(boxOfficeStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат кассовых сборов"})
			return
		}
		movie.BoxOffice = boxOffice
	}

	if _, err := c.FormFile("movieURL"); err == nil {
		prePath := movie.MovieURL
		filePath, err := utils.ReplaceFile(c, prePath, "movieURL", name, "media/movie-videos")
		if err != nil {
			log.Println("Ошибка получения пути файла: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintln("Ошибка получения пути файла: ", err)})
			return
		}
		movie.MovieURL = fmt.Sprintf("http://192.168.1.9:8080/%s", filePath)
	}
	if _, err := c.FormFile("posterURL"); err == nil {
		prePath := movie.PosterURL
		filePath, err := utils.ReplaceFile(c, prePath, "posterURL", name, "media/movie-posters")
		if err != nil {
			log.Println("Ошибка получения пути файла: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintln("Ошибка получения пути файла: ", err)})
			return
		}
		movie.PosterURL = fmt.Sprintf("http://192.168.1.9:8080/%s", filePath)
	}
	if _, err := c.FormFile("trailerURL"); err == nil {
		prePath := movie.TrailerURL
		filePath, err := utils.ReplaceFile(c, prePath, "trailerURL", name, "media/movie-trailers")
		if err != nil {
			log.Println("Ошибка получения пути файла: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintln("Ошибка получения пути файла: ", err)})
			return
		}
		movie.TrailerURL = fmt.Sprintf("http://192.168.1.9:8080/%s", filePath)
	}

	if genreStr := c.PostForm("genres"); genreStr != "" {
		ids := strings.Split(genreStr, ",")
		var newGenres []models.Genre
		for _, idStr := range ids {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				log.Printf("Неверный ID категории: %s", idStr)
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Неверный ID категории: %s", idStr)})
				return
			}
			var cat models.Genre
			if err := utils.Db.First(&cat, id).Error; err != nil {
				log.Printf("Категория %d не найдена: %v", id, err)
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Категория %d не найдена", id)})
				return
			}
			newGenres = append(newGenres, cat)
		}
		if err := utils.Db.Model(&movie).Association("Genres").Replace(newGenres); err != nil {
			log.Printf("Ошибка обновления категорий: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления категорий"})
			return
		}
	}

	if movieMembersStr := c.PostForm("movieMembers"); movieMembersStr != "" {
		var movieMembers []models.MovieMember
		if err := json.Unmarshal([]byte(movieMembersStr), &movieMembers); err != nil {
			log.Println("Ошибка парсинга участников фиильмов: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintln("Ошибка парсинга участников фиильмов: ", err)})
			return
		}
		if err := utils.Db.Model(&movie).Association("MovieMembers").Replace(movieMembers); err != nil {
			log.Printf("Ошибка обновления участников: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления участников"})
			return
		}
	}

	if err := utils.Db.Save(&movie).Error; err != nil {
		log.Printf("Ошибка обновления фильма: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Ошибка обновления фильма: %v", err)})
		return
	}

	if err := utils.Db.Preload("Genres").Preload("MovieMembers.Member").First(&movie, id).Error; err != nil {
		log.Println("Ошиибка поулчения фильма: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintln("Ошиибка поулчения фильма: ", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Фильм успешно обновлен", "movie": movie})
}
