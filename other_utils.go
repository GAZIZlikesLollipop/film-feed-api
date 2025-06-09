package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func getCategories(c *gin.Context) {
	var categories []Category
	if err := db.Find(&categories).Error; err != nil {
		log.Println("Ошибка поулчения категорий", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения категорий"})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func deleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&Category{}, id).Error; err != nil {
		log.Println("Ошибка удаления категории", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления категорий"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Категория успешно удалена"})
}

func addCategory(c *gin.Context) {
	var category Category
	if err := c.ShouldBindJSON(&category); err != nil {
		log.Println("Введены не верные данные", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Введены не верные данные"})
		return
	}
	if err := db.Create(&category).Error; err != nil {
		log.Println("Ошибка добавленя категории", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка добавелния категории"})
		return
	}
	c.JSON(http.StatusCreated, category)
}

func updateCategory(c *gin.Context) {
	id := c.Param("id")
	var category Category
	if err := db.First(&category, id).Error; err != nil {
		log.Println("Категория не найдена", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Категория не найдена"})
		return
	}
	var newCategory Category
	if err := c.ShouldBindJSON(&newCategory); err != nil {
		log.Printf("Неверные данные: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}
	category.Name = newCategory.Name
	if err := db.Save(&category).Error; err != nil {
		log.Println("Ошибка обновления категории: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления категории"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "категория успешно обновлена"})
}
