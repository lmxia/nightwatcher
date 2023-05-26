package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/lmxia/nightwatcher/api/v1/gaia"
)

func addGaiaRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/gaia")
	router.GET("/clusters", gaia.GetClusters)
	router.GET("/clusters/:namespace/:cluster/:label", gaia.GetClusterLabel)
}
