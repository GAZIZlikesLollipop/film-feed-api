package handlers

import (
	"log"
	"net/http"

	"api/internal/models"
	utils "api/pkg"

	"github.com/gin-gonic/gin"
)

func GetGenres(c *gin.Context) {
	var categories []models.Genre
	if err := utils.Db.Find(&categories).Error; err != nil {
		log.Println("Ошибка поулчения категорий", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения категорий"})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func DeleteGenre(c *gin.Context) {
	id := c.Param("id")
	if err := utils.Db.Delete(&models.Genre{}, id).Error; err != nil {
		log.Println("Ошибка удаления категории", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления категорий"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Категория успешно удалена"})
}

func AddGenre(c *gin.Context) {
	var category models.Genre
	if err := c.ShouldBindJSON(&category); err != nil {
		log.Println("Введены не верные данные", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Введены не верные данные"})
		return
	}
	if err := utils.Db.Create(&category).Error; err != nil {
		log.Println("Ошибка добавленя категории", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка добавелния категории"})
		return
	}
	c.JSON(http.StatusCreated, category)
}

func UpdateGenre(c *gin.Context) {
	id := c.Param("id")
	var category models.Genre
	if err := utils.Db.First(&category, id).Error; err != nil {
		log.Println("Категория не найдена", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Категория не найдена"})
		return
	}
	var newCategory models.Genre
	if err := c.ShouldBind(&newCategory); err != nil {
		log.Printf("Неверные данные: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}
	category.Name = newCategory.Name
	if err := utils.Db.Save(&category).Error; err != nil {
		log.Println("Ошибка обновления категории: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления категории"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "категория успешно обновлена"})
}
