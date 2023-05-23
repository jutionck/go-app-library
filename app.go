package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Product struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost,
		dbPort,
		dbUser,
		dbPassword,
		dbName,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Product{}, &User{})
	db.Create(&User{
		Id:       "001",
		Username: "admin",
		Password: "password",
	})

	r := gin.Default()

	r.GET("/products", func(ctx *gin.Context) {
		var products []Product
		if err := db.Find(&products).Error; err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": products,
		})
	})

	r.GET("/products/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var product Product
		if err := db.First(&product, "id=?", id).Error; err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": product,
		})
	})

	r.POST("/products", func(ctx *gin.Context) {
		var payload Product
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": err.Error(),
			})
			return
		}

		if err := db.Save(&payload).Error; err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"code": http.StatusCreated,
			"data": payload,
		})
	})

	r.PUT("/products", func(ctx *gin.Context) {
		var payload Product
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": err.Error(),
			})
			return
		}

		if err := db.Save(&payload).Error; err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"code": http.StatusCreated,
			"data": payload,
		})
	})

	r.DELETE("/products/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		if err := db.Delete(&Product{}, "id=?", id).Error; err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		ctx.String(http.StatusNoContent, "")
	})

	r.POST("/auth", func(ctx *gin.Context) {
		var payload User
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		if err := db.First(&payload, "username=?", payload.Username).Error; err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"code":  http.StatusCreated,
			"token": "initoken",
		})
	})

	r.POST("/auth/logout", func(ctx *gin.Context) {})

	apiHost := os.Getenv("API_HOST")
	apiPort := os.Getenv("API_PORT")
	host := fmt.Sprintf("%s:%s", apiHost, apiPort)
	err = r.Run(host)
	if err != nil {
		panic(err)
	}

}
