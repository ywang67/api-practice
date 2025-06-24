package router

import (
	"api-project/restful-api/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) *gin.Engine {
	api := r.Group("/api/v1")
	{
		cm := api.Group("/cablemodems")
		{
			cm.GET("/by-mac", handler.CableModemsByMac)
			cm.GET("/by-cmts", handler.CableModemsByCmts)
			cm.GET("/by-poller", handler.CableModemsByPoller)
		}
	}

	return r
}
