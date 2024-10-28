package rest

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"

	"github.com/mfsyahrz/image_feed_api/internal/common/file"
	"github.com/mfsyahrz/image_feed_api/internal/common/pagination"
	"github.com/mfsyahrz/image_feed_api/internal/common/util"
	"github.com/mfsyahrz/image_feed_api/internal/domain/repository"
	"github.com/mfsyahrz/image_feed_api/internal/service/post"
)

type PostHandler struct {
	service post.PostService
}

func NewPostHandler(service post.PostService) *PostHandler {
	return &PostHandler{service}
}

func (h *PostHandler) Find(c echo.Context) error {
	var (
		resp = new(PostResponse)
		req  GetPaginatedPostsRequest
		err  error
	)

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

	cursor, err := pagination.DecodePostCursor(util.StringVal(req.Cursor))
	if err != nil {
		return err
	}

	posts, nextCursor, err := h.service.GetPaginated(
		c.Request().Context(),
		repository.GetPostInput{
			Cursor: cursor,
			Limit:  req.Limit,
		})
	if err != nil {
		return err
	}
	nextCursorStr := nextCursor.Encode()

	resp = NewPostsResponse(posts, &nextCursorStr, "Posts Retrieved Succesfully")

	return nil
}

func (h *PostHandler) Create(c echo.Context) error {
	var (
		resp = new(PostResponse)
		req  post.CreatePostInput
		err  error
	)

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

	req.SrcImgPath, req.DisplayImgPath, err = h.saveImage(c)
	if err != nil {
		return err
	}

	post, err := h.service.Save(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	posts, _, err := h.service.GetPaginated(c.Request().Context(), repository.GetPostInput{
		IDs: []int64{post.ID},
	})
	if err != nil {
		return err
	}

	resp = NewPostsResponse(posts, nil, "Posts Created Succesfully")

	return nil
}

func (h *PostHandler) saveImage(c echo.Context) (string, string, error) {
	srcFile, err := c.FormFile("image")
	if err != nil {
		return "", "", fmt.Errorf("image upload error: %+v", err)
	}

	fh := file.NewFileHeader(srcFile)
	return h.service.SaveImage(c.Request().Context(), fh)
}
