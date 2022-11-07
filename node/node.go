package node

import (
	"database/sql"
	"gin_backend/infra"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

var db = infra.GetDB()
var err error

type nodes struct {
	gorm.Model
	City      string
	Name      string
	Ename     string
	GpsLati   string
	GpsLong   string
	IsDeleted sql.NullBool `gorm:"default:false"`
}

func Config(api *gin.RouterGroup) {
	_ = db.DB.AutoMigrate(&nodes{})
	api.GET("/nodes", getNodes)
	api.GET("/nodes/:id", getNode)
}

func getNodes(c *gin.Context) {
	var n []nodes
	if err = db.DB.Find(&n).Error; err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		log.Fatalf(err.Error())
	} else {
		c.JSON(http.StatusOK, n)
	}
}

func getNode(c *gin.Context) {
	id := c.Param("id")
	var n nodes
	if err = db.DB.Where("id = ?", id).First(&n).Error; err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
	} else {
		c.JSON(http.StatusOK, n)
	}
}
