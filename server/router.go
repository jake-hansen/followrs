package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jake-hansen/followrs/handlers"
	"github.com/jake-hansen/followrs/repositories"
	"github.com/jake-hansen/followrs/services"
)

// NewRouter returns a router configured with handlers for configured
// endpoints.
func NewRouter(env string) *gin.Engine {
	setGinEnvironment(env)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group("v1")
	handlers.NewHealthHandler(v1, services.NewSimpleHealthService(repositories.NewSimpleHealthRepository()))

	return router
}

func setGinEnvironment(env string) {
	if env == "prod" {
		gin.SetMode("release")
	}
}
