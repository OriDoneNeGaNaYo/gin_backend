package main

import (
	"database/sql"
	"gin_backend/infra"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

// gin-swagger middleware
// swagger embed files

var con = getInfra()

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

func getInfra() infra.Database {
	infra.LoadEnv()
	return infra.NewDatabase()
}

/* 아래 항목이 swagger에 의해 문서화 된다. */
// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /api/v1
func main() {
	_ = con.DB.AutoMigrate(&nodes{})
	router := gin.Default()
	v1 := router.Group("/api/v1")
	v1.GET("/nodes", getNodes)
	v1.GET("/nodes/:id", getNode)
	router.Run(":8080")
}

func getNodes(c *gin.Context) {
	var n []nodes
	if err = con.DB.Find(&n).Error; err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		log.Fatalf(err.Error())
	} else {
		c.JSON(http.StatusOK, n)
	}
}

func getNode(c *gin.Context) {
	id := c.Param("id")
	var n nodes
	if err = con.DB.Where("id = ?", id).First(&n).Error; err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
	} else {
		c.JSON(http.StatusOK, n)
	}
}

func postNode(c *gin.Context) {

}
