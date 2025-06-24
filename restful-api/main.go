package main

import (
	"api-project/pkg/dbservice"
	"api-project/restful-api/router"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	dbService := dbservice.DbService
	r.Use(func(c *gin.Context) {
		c.Set("dbRead", dbService.DbReader)
		c.Set("dbWrite", dbService.DbWriter)
		c.Next()
	})
	router.SetupRouter(r)
	r.Run(":8080")
}
