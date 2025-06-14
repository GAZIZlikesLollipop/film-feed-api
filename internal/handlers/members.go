package handlers

import (
	"api/internal/models"
	utils "api/pkg"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func AddMember(c *gin.Context) {
	name := c.PostForm("name")
	roles := c.PostFormArray("roles")
	layout := "2006-01-02"
	birthDate, err := time.Parse(layout, c.PostForm("birthDate"))
	if err != nil {
		log.Println("Ошибка преобразования даты: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintln("Ошибка преобразования даты: ", err)})
		return
	}
	biography := c.PostForm("biography")
	nationality := c.PostForm("nationality")

	featuredFilmsStr := c.PostForm("featuredFilms")
	var featuredFilms []models.Movie
	if featuredFilmsStr != "" {
		ids := strings.Split(featuredFilmsStr, ",")
		for _, idStr := range ids {
			movieID, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				log.Printf("Неверный ID участника: %s", idStr)
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Неверный ID участника: %s", idStr)})
				return
			}
			var movie models.Movie
			if err := utils.Db.Preload("Genres").Preload("MovieMembers.Member").First(&movie, movieID).Error; err != nil {
				log.Printf("Фильм %d не найдена: %v", movieID, err)
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Фильм %d не найдена", movieID)})
				return
			}
			featuredFilms = append(featuredFilms, movie)
		}
		if len(featuredFilms) <= 0 {
			featuredFilms = []models.Movie{}
		}
	} else {
		featuredFilms = []models.Movie{}
	}

	photoPath, err := utils.SaveFile(c, "photo", name, "member-photos")
	if err != nil {
		log.Println("Ошибка получения пути файла: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintln("Ошибка получения пути файла: ", err)})
		return
	}

	member := models.Member{
		Name:          name,
		Roles:         roles,
		BirthDate:     birthDate,
		Biography:     biography,
		Nationality:   nationality,
		Photo:         fmt.Sprintf("http://192.168.1.9:8080/%s", photoPath),
		FeaturedFilms: featuredFilms,
	}

	if deathDateStr := c.PostForm("deathDate"); deathDateStr != "" {
		if deathDate, err := time.Parse(layout, deathDateStr); deathDateStr != "" {
			if err != nil {
				member.DeathDate = nil
			} else {
				member.DeathDate = &deathDate
			}
		}
	}

	if err := utils.Db.Create(&member).Error; err != nil {
		log.Println("Ошибка добовления учатсниика: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintln("Ошибка добовления учатсниика: ", err)})
		return
	}

	c.JSON(http.StatusCreated, member)
}

func GetMembers(c *gin.Context) {
	var members []models.Member
	if err := utils.Db.Preload("FeaturedFilms").Find(&members).Error; err != nil {
		log.Println("Ошибка получения учатсников: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintln("Ошибка получения участников: ", err)})
		return
	}
	c.JSON(http.StatusOK, members)
}

func GetMember(c *gin.Context) {
	id := c.Param("id")
	var member models.Member
	if err := utils.Db.Preload("FeaturedFilms").First(&member, id).Error; err != nil {
		log.Println("Ошибка получения учатсника: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintln("Ошибка получения участника: ", err)})
		return
	}
	c.JSON(http.StatusOK, member)
}

func DeleteMember(c *gin.Context) {
	id := c.Param("id")
	var member models.Member

	if err := utils.Db.Preload("FeaturedFilms").First(&member, id).Error; err != nil {
		log.Println("Ошибка получения учатника", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получениия участника"})
		return
	}

	if member.Photo != "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Println("Ошибка получшения пользовтельской диреткории: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintln("Ошибка получшения пользовтельской диреткории: ", err)})
			return
		}
		file, err := url.Parse(member.Photo)
		if err != nil {
			log.Println("Ошибка прасиинга Url", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintln("Ошиибка парсиинга URL", err)})
			return
		}
		if err := os.Remove(filepath.Join(homeDir, file.Path)); err != nil {
			log.Println("Ошибка удаления файла: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintln("Ошибка удаления файла: ", err)})
			return
		}
	}

	if err := utils.Db.Model(&member).Association("FeaturedFilms").Clear(); err != nil {
		log.Printf("Ошибка очистки связей фильмов: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintln("Ошибка очистки связей фильмов: ", err)})
		return
	}

	if err := utils.Db.Delete(&member).Error; err != nil {
		log.Println("Ошибка удаления учатника", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления участника"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Участник успшено удален"})
}

func UpdateMember(c *gin.Context) {
	var member models.Member
	id := c.Param("id")

	if err := utils.Db.Preload("FeaturedFilms").First(&member, id).Error; err != nil {
		log.Println("Фильм не найден", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Фильм не найден"})
		return
	}

	if name := c.PostForm("name"); name != "" {
		member.Name = name
	}
	if roles := c.PostFormArray("roles"); len(roles) != 0 {
		member.Roles = roles
	}
	layout := "2006-01-02"
	birthDateStr := c.PostForm("birthDate")
	if birthDate, err := time.Parse(layout, birthDateStr); birthDateStr != "" {
		if err != nil {
			log.Println("Ошибка преобразования даты: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintln("Ошибка преобразования даты: ", err)})
			return
		}
		member.BirthDate = birthDate
	}

	if deathDateStr := c.PostForm("deathDate"); deathDateStr != "" {
		if deathDate, err := time.Parse(layout, deathDateStr); deathDateStr != "" {
			if err != nil {
				member.DeathDate = nil
			} else {
				member.DeathDate = &deathDate
			}
		}
	}

	if biography := c.PostForm("biography"); biography != "" {
		member.Biography = biography
	}
	if nationality := c.PostForm("nationality"); nationality != "" {
		member.Nationality = nationality
	}

	if featuredFilmsStr := c.PostForm("featuredFilms"); featuredFilmsStr != "" {
		var featuredFilms []models.Movie
		ids := strings.Split(featuredFilmsStr, ",")
		for _, idStr := range ids {
			movieID, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				log.Printf("Неверный ID участника: %s", idStr)
				// c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Неверный ID участника: %s", idStr)})
				// return
				break
			}
			var movie models.Movie
			if err := utils.Db.Preload("Genres").Preload("MovieMembers.Member").First(&movie, movieID).Error; err != nil {
				log.Printf("Фильм %d не найдена: %v", movieID, err)
				// c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Фильм %d не найдена", movieID)})
				// return
				break
			}
			featuredFilms = append(featuredFilms, movie)
		}
		if len(featuredFilms) <= 0 {
			featuredFilms = []models.Movie{}
		}
		if err := utils.Db.Model(&member).Association("FeaturedFilms").Replace(featuredFilms); err != nil {
			log.Printf("Ошибка обновления фильмов: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления фильмов"})
			return
		}
		member.FeaturedFilms = featuredFilms
	}

	if _, err := c.FormFile("photo"); err == nil {
		prePath := member.Photo
		filePath, err := utils.ReplaceFile(c, prePath, "photo", member.Name, "member-photos")
		if err != nil {
			log.Println("Ошибка получения пути файла: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintln("Ошибка получения пути файла: ", err)})
			return
		}
		member.Photo = fmt.Sprintf("http://192.168.1.9:8080/%s", filePath)
	}

	if err := utils.Db.Save(&member).Error; err != nil {
		log.Printf("Ошибка обновления учатнсика: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Ошибка обновления участника: %v", err)})
		return
	}

	if err := utils.Db.Preload("FeaturedFilms").First(&member, id).Error; err != nil {
		log.Println("Фильм не найден", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Фильм не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": member})
}
