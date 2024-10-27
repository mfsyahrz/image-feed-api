package rest

import "github.com/mfsyahrz/image_feed_api/internal/domain/entity"

const (
	ResponseSuccess = "Success"
	IDEmpty         = "ID Cannot Be Empty"
)

type BaseResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PostResponse struct {
	BaseResponse
	Cursor *string `json:"cursor,omitempty"`
}

func NewPostsResponse(data entity.Posts, cursor *string, msg string) *PostResponse {
	return &PostResponse{
		BaseResponse: BaseResponse{
			Success: true,
			Message: msg,
			Data:    data,
		},
		Cursor: cursor,
	}
}

type CommentResponse struct {
	BaseResponse
}

func NewCommentsResponse(data entity.Comments, msg string) *CommentResponse {
	return &CommentResponse{
		BaseResponse: BaseResponse{
			Success: true,
			Message: msg,
			Data:    data,
		},
	}
}

type GetPaginatedPostsRequest struct {
	Cursor *string 	  `json:"cursor"`
	Limit  int     	  `json:"limit"`
	CommentLimit int  `json:"comment_limit"`
}

type DeleteCommentRequest struct {
	ID     int64
	PostID int64
}
