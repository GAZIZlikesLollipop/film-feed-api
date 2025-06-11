package utils

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var Db *gorm.DB

func SaveFile(
	c *gin.Context,
	field string,
	name string,
	directory string,
) (string, error) {
	file, err := c.FormFile(field)
	if err != nil {
		log.Println("Ошиибка получения файла: ", err)
		return "", err
	}

	if err := os.MkdirAll(directory, 0755); err != nil {
		log.Println("Ошибка создания директории: ", err)
		return "", err
	}

	fileName := fmt.Sprintf("%s-%s%s", name, uuid.New(), strings.ToLower(filepath.Ext(file.Filename)))
	filePath := filepath.Join(directory, fileName)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		log.Printf("Ошибка сохранения файла: %v", err)
		return "", err
	}

	return filePath, nil
}

func ReplaceFile(
	c *gin.Context,
	prePath string,
	field string,
	name string,
	directory string,
) (string, error) {
	absolutePath := "/home/lollipop/dev/film-feed"
	preFile, err := url.Parse(prePath)
	if err != nil {
		log.Println("Ошибка прасиинга Url", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintln("Ошиибка парсиинга URL", err)})
		return "", err
	}
	if err := os.Remove(fmt.Sprintf("%s%s", absolutePath, preFile.Path)); err != nil {
		log.Println("Ошиибка удаления файла: ", err)
		return "", err
	}

	file, err := c.FormFile(field)
	if err != nil {
		log.Println("Ошиибка получения файла: ", err)
		return "", err
	}

	if err := os.MkdirAll(directory, 0755); err != nil {
		log.Println("Ошибка создания директории: ", err)
		return "", err
	}

	fileName := fmt.Sprintf("%s-%s%s", name, uuid.New(), strings.ToLower(filepath.Ext(file.Filename)))
	filePath := filepath.Join(directory, fileName)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		log.Printf("Ошибка сохранения файла: %v", err)
		return "", err
	}

	return filePath, nil
}
