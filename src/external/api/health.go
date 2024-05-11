package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcsouza/fiap-tech-fast-food/src/external/api/infra/config"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func RegisterHealthRoutes(gRouter *gin.Engine) {
	gRouter.GET("/health/liveness", getLivenessHandler)
	gRouter.GET("/health/readiness", getReadinessHandler)
}

// Liveness godoc
// @Summary Liveness probe
// @Description Liveness probe
// @Tags Health Routes
// @Produce  json
// @Success 200 {object} map[string]string
// @Router /health/liveness [get]
func getLivenessHandler(c *gin.Context) {
	c.JSON(healthCheck())
}

// Readiness godoc
// @Summary Readiness probe
// @Description Readiness probe
// @Tags Health Routes
// @Produce  json
// @Success 200 {object} map[string]string
// @Router /health/readiness [get]
func getReadinessHandler(c *gin.Context) {
	c.JSON(healthCheck())
}

func healthCheck() (code int, obj any) {
	if databaseHealth() {
		return http.StatusOK, struct{ Status string }{Status: "OK"}
	} else {
		return http.StatusBadGateway, struct{ Status string }{Status: "FAILED"}
	}
}

func databaseHealth() (healthStatus bool) {
	healthStatus = false
	client, err := config.GetMongoClient()
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return
	}
	healthStatus = true
	return
}
