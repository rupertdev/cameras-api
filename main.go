package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

type CreateCameraInput struct {
	CameraModel string `json:"model" binding:"required"`
	CameraFormat string `json:"format" binding:"required"`
	CameraManufacturer string `json:"manufacturer" binding:"required"`
}

type Camera struct {
	gorm.Model
	Name string
	Manufacturer string
	Format string
}

func main() {
	db_url := os.Getenv("DATABASE_URL")
db, err := gorm.Open(postgres.Open(&db_url), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&Camera{})

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.POST("/cameras", func(c *gin.Context) {
		var input CreateCameraInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		camera := Camera{Name: input.CameraModel, Manufacturer: input.CameraManufacturer, Format: input.CameraFormat}
		db.Create(&camera)

	})

	router.Run(":" + port)
}
