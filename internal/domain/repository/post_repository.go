package repository

import (
	"context"

	"github.com/mfsyahrz/image_feed_api/internal/common/pagination"
	"github.com/mfsyahrz/image_feed_api/internal/domain/entity"
)

type GetPostInput struct {
	Cursor *pagination.PostCursor
	Limit  int 		 
	CommentLimit int
	IDs    []int64
}

type PostRepository interface {
	FetchPaginated(ctx context.Context, input GetPostInput) (entity.Posts, *pagination.PostCursor, error)
	FetchOne(ctx context.Context, input GetPostInput) (*entity.Post, error)
	Save(ctx context.Context, post *entity.Post) error
}
