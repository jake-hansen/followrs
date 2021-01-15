package server

import (
	"github.com/jake-hansen/followrs/config"
	"github.com/jake-hansen/followrs/domain"
	"github.com/jake-hansen/followrs/repositories/apis/twitter"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jake-hansen/followrs/handlers"
	"github.com/jake-hansen/followrs/middleware"
	"github.com/jake-hansen/followrs/repositories"
	"github.com/jake-hansen/followrs/services"
)

// NewRouter returns a router configured with handlers for configured
// endpoints.
func NewRouter(env string, startTime time.Time) *gin.Engine {
	setGinEnvironment(env)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.PublicErrorHandler())

	v1 := router.Group("v1")
	handlers.NewHealthHandler(v1, services.NewSimpleHealthService(repositories.NewSimpleHealthRepository(startTime)))

	handlers.NewUsersHandler(v1, createTwitterService())

	return router
}

func setGinEnvironment(env string) {
	if env == "prod" {
		gin.SetMode("release")
	}
}

func createTwitterService() *domain.TwitterService {
	config := config.GetConfig()
	apiKey := config.GetString("secrets.twitter.api.key")
	apiSecretKey := config.GetString("secrets.twitter.api.secret")
	apiBearerToken := config.GetString("secrets.twitter.api.bearer")

	twitterRepo, _ := twitter.NewTwitterAPI("https://api.twitter.com/2", apiKey, apiSecretKey, apiBearerToken)
	repoPtr := domain.TwitterRepository(twitterRepo)

	service := services.NewTwitterService(&repoPtr)

	return &service
}
