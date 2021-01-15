package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jake-hansen/followrs/apperrors"
	"github.com/jake-hansen/followrs/domain"
	"net/http"
	"strings"
)

type UsersHandler struct {
	TwitterService *domain.TwitterService
}

func NewUsersHandler(parentGroup *gin.RouterGroup, twitterService *domain.TwitterService) {
	handler := &UsersHandler{
		TwitterService: twitterService,
	}

	usersGroup := parentGroup.Group("users")
	{
		twitterGroup := usersGroup.Group("/twitter")
		{
			twitterGroup.GET("/:username", func(c *gin.Context) {
				username := c.Param("username")
				handler.GetTwitterUser(username, c)
			})
		}
	}
}

func (u *UsersHandler) GetTwitterUser(username string, c *gin.Context) {
	user, err := (*u.TwitterService).GetUser(username)

	if err == nil {
		c.JSON(http.StatusOK, *user)
	} else {
		apiError := err

		if strings.Contains(err.Error(), "user not found") {
			apiError = &apperrors.APIError{
				Status:  http.StatusNotFound,
				Err:     err,
				Message: fmt.Sprintf("the user [%s] was not found", username),
			}
		}
		c.Error(apiError).SetType(gin.ErrorTypePublic)
	}
}