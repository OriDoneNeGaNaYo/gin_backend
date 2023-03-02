package main

import (
	busStop "gin_backend/bus-stop"
	_ "gin_backend/docs"
	"gin_backend/user"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

/* 아래 항목이 swagger에 의해 문서화 된다. */
// @title 우리동네가나요 API
// @version 1.0
// @description This is a sample server OriDoneNeGaNaYo server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @host OriDoneNeGaNaYo.swagger.io
// @BasePath /api/v1
func main() {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1 := router.Group("/api/v1")
	busStop.Config(v1)
	user.Config(v1)

	_ = router.Run(":8080")
}
