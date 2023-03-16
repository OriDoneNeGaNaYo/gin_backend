package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	infra "main-server/infra"
	"net/http"
)

var db = infra.GetDB()
var writer infra.Writer
var err error

type user struct {
	gorm.Model
	ID   string `gorm:"primaryKey"`
	Name string `gorm:"uniqueKey"`
}

func Config(api *gin.RouterGroup) {
	_ = db.DB.AutoMigrate(&user{})
	api.GET("/users", getAllUser)
	api.POST("/users", postUser)
	api.GET("/users/:id", getUser)
}

type userRequest struct {
	ID   string `form:"id" binding:require`
	Name string `form:"name" binding:require`
}

func postUser(c *gin.Context) {
	var req userRequest
	if err := c.ShouldBind(&req); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	newUser := user{
		ID:   req.ID,
		Name: req.Name,
	}

	if err := db.DB.Create(&newUser).Error; err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

func getAllUser(c *gin.Context) {
	var data []user

	if err = db.DB.Find(&data).Error; err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		log.Fatalf(err.Error())
	} else {
		//getKafka()
		c.JSON(http.StatusOK, data)
	}
}

func getUser(c *gin.Context) {
	id := c.Param("id")
	var data user
	if err = db.DB.Where("id = ?", id).First(&data).Error; err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}
