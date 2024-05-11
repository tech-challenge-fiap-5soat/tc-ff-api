package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterSwaggerRoutes(gRouter *gin.Engine) {
	gRouter.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
