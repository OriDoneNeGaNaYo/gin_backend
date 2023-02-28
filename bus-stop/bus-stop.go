package bus_stop

import (
	"fmt"
	"gin_backend/infra"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

var db = infra.GetDB()
var writer infra.Writer
var err error

// id,city,name,gps_lati,gps_long,collected_at
type busStop struct {
	gorm.Model
	ID          string `gorm:"primaryKey"`
	City        string `gorm:"index"`
	Name        string
	GpsLati     string
	GpsLong     string
	CollectedAt string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	updatedAt   time.Time `gorm:"autoUpdateTime"`
}

func Config(api *gin.RouterGroup) {
	_ = db.DB.AutoMigrate(&busStop{})
	api.GET("/bus-stops", getAllBusStops)
	api.GET("/bus-stops/:id", getBusStop)
	//go writer.GetKafkaWriter("9091", "bus", 1)
}

type PaginationQuery struct {
	Key    string `form:"key" binding:"omitempty"`
	Page   int    `form:"page" binding:"omitempty,min=0"`
	Size   int    `form:"size" binding:"omitempty,min=1,max=100"`
	SortBy string `form:"sort_by" binding:"omitempty,oneof=id city name gps_lati gps_long"`
	Sort   string `form:"sort" binding:"omitempty,oneof=asc desc"`
}

func getAllBusStops(c *gin.Context) {
	var data []busStop
	var paginationQuery PaginationQuery

	if err := c.BindQuery(&paginationQuery); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if paginationQuery.Sort == "" {
		paginationQuery.Sort = "asc"
	}
	if paginationQuery.SortBy == "" {
		paginationQuery.SortBy = "name"
	}

	pagination := infra.Pagination{
		Limit: paginationQuery.Size,
		Page:  paginationQuery.Page,
		Sort:  fmt.Sprintf("%s %s", paginationQuery.SortBy, paginationQuery.Sort),
	}

	// 키가 존재할 경우 정류장 이름과 도시 중 하나라도 포함 되는 정류장을 가져온다.
	chainQuery := db.DB
	if paginationQuery.Key != "" {
		chainQuery = chainQuery.Or("name LIKE ?", "%"+paginationQuery.Key+"%")
		chainQuery = chainQuery.Or("city LIKE ?", "%"+paginationQuery.Key+"%")
	}

	if err = chainQuery.Scopes(infra.Paginate(&data, &pagination, chainQuery)).Find(&data).Error; err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		log.Fatalf(err.Error())
	} else {
		//getKafka()
		c.JSON(http.StatusOK, data)
	}
}

func getBusStop(c *gin.Context) {
	id := c.Param("id")
	var n busStop
	if err = db.DB.Where("id = ?", id).First(&n).Error; err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
	} else {
		c.JSON(http.StatusOK, n)
	}
}
