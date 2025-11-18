package web

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupRoutes(db *gorm.DB) *gin.Engine {
	DB = db

	r := gin.Default()
	r.LoadHTMLGlob("pkg/web/templates/*.html")

	r.GET("/", showAllLogs) //first time loading
	r.POST("/", filterLogs) // after searching

	return r
}
