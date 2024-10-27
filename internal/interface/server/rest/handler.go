package rest

import (
	"github.com/labstack/echo"
	"github.com/mfsyahrz/image_feed_api/internal/interface/ioc"
)

type Handler struct {
	PostHandler    *PostHandler
	CommentHandler *CommentHandler
}

func SetupHandler(artifact *ioc.IOC) *Handler {
	return &Handler{
		PostHandler:    NewPostHandler(artifact.PostService),
		CommentHandler: NewCommentHandler(artifact.CommentService),
	}
}

func bindRequest(c echo.Context, req interface{}) error {
	if err := c.Bind(req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	return nil
}
