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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Println("Ошибка получшения пользовтельской диреткории: ", err)
		return "", err
	}

	absolutePath := filepath.Join(homeDir, "media", directory)

	file, err := c.FormFile(field)
	if err != nil {
		log.Println("Ошиибка получения файла: ", err)
		return "", err
	}

	if err := os.MkdirAll(absolutePath, 0755); err != nil {
		log.Println("Ошибка создания директории: ", err)
		return "", err
	}

	fileName := fmt.Sprintf("%s-%s%s", name, uuid.New(), strings.ToLower(filepath.Ext(file.Filename)))
	saveFilePath := filepath.Join(absolutePath, fileName)

	if err := c.SaveUploadedFile(file, saveFilePath); err != nil {
		log.Printf("Ошибка сохранения файла: %v", err)
		return "", err
	}

	filePath := filepath.Join("media", directory, fileName)

	return filePath, nil
}

func ReplaceFile(
	c *gin.Context,
	prePath string,
	field string,
	name string,
	directory string,
) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Println("Ошибка получшения рабочей диреткории: ", err)
		return "", err
	}
	preFile, err := url.Parse(prePath)
	if err != nil {
		log.Println("Ошибка прасиинга Url", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintln("Ошиибка парсиинга URL", err)})
		return "", err
	}
	delPath := filepath.Join(homeDir, preFile.Path)
	_, exist := os.Stat(delPath)
	if exist == nil {
		if err := os.Remove(delPath); err != nil {
			log.Println("Ошиибка удаления файла: ", err)
			return "", err
		}
	}

	file, err := c.FormFile(field)
	if err != nil {
		log.Println("Ошиибка получения файла: ", err)
		return "", err
	}

	absolutePath := filepath.Join(homeDir, "media", directory)

	if err := os.MkdirAll(absolutePath, 0755); err != nil {
		log.Println("Ошибка создания директории: ", err)
		return "", err
	}

	fileName := fmt.Sprintf("%s-%s%s", name, uuid.New(), strings.ToLower(filepath.Ext(file.Filename)))
	saveFilePath := filepath.Join(absolutePath, fileName)

	if err := c.SaveUploadedFile(file, saveFilePath); err != nil {
		log.Printf("Ошибка сохранения файла: %v", err)
		return "", err
	}

	filePath := filepath.Join("media", directory, fileName)

	return filePath, nil
}
