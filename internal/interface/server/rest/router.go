package rest

import (
	"net/http"

	"github.com/labstack/echo"
)

func SetupRouter(server *echo.Echo, handler *Handler) {

	// - health check
	server.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "service up and running... ")
	})

	postRoute := server.Group("/posts")
	{
		postRoute.POST("", handler.PostHandler.Create)
		postRoute.GET("", handler.PostHandler.Find)
	}

	server.Static("/storage/images", "storage/images")

	commentRoute := postRoute.Group("/:postID/comments")
	{
		commentRoute.POST("", handler.CommentHandler.Save)
		commentRoute.DELETE("/:id", handler.CommentHandler.Delete)

	}

}
