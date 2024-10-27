package repository

import (
	"context"

	"github.com/mfsyahrz/image_feed_api/internal/domain/entity"
)

type GetCommentInput struct {
	PostIDs []int64
}

type CommentRepository interface {
	Fetch(ctx context.Context, input GetCommentInput) (entity.Comments, error)
	Save(ctx context.Context, comment *entity.Comment) error
	Delete(ctx context.Context, comment *entity.Comment) error 
}
