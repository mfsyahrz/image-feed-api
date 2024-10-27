package rest

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/spf13/cast"

	"github.com/mfsyahrz/image_feed_api/internal/domain/entity"
	"github.com/mfsyahrz/image_feed_api/internal/service/comment"
)

type CommentHandler struct {
	service comment.CommentService
}

func NewCommentHandler(service comment.CommentService) *CommentHandler {
	return &CommentHandler{service}
}

func (h *CommentHandler) Save(c echo.Context) error {
	var (
		resp = new(CommentResponse)
		req  entity.Comment
		err  error
	)

	req.PostID = cast.ToInt64(c.Param("postID"))

	if err := bindRequest(c, &req); err != nil {
		resp.Message = err.Error()
		return c.JSON(http.StatusBadRequest, resp)
	}

	defer func() {
		if err != nil {
			resp.Message = err.Error()
		}
		c.JSON(http.StatusCreated, resp)
	}()

	err = h.service.Save(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	resp = NewCommentsResponse(entity.Comments{&req}, "Comment Saved Succesfully")

	return nil
}

func (h *CommentHandler) Delete(c echo.Context) error {
	var (
		resp = new(CommentResponse)
		req  DeleteCommentRequest
		err  error
	)

	req.ID = cast.ToInt64(c.Param("id"))
	req.PostID = cast.ToInt64(c.Param("postID"))

	if err := bindRequest(c, &req); err != nil {
		resp.Message = err.Error()
		return c.JSON(http.StatusBadRequest, resp)
	}

	defer func() {
		if err != nil {
			resp.Message = err.Error()
		}
		c.JSON(http.StatusOK, resp)
	}()

	err = h.service.Delete(c.Request().Context(), &entity.Comment{
		ID:     req.ID,
		PostID: req.PostID,
	})
	if err != nil {
		return err
	}

	resp = NewCommentsResponse(nil, fmt.Sprintf("CommentID %d Deleted Succesfully", req.ID))

	return nil
}
