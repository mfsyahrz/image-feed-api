package repository

import (
	"context"

	"github.com/mfsyahrz/image_feed_api/internal/domain/entity"
)

type GetCommentInput struct {
	PostIDs []int64
}

type CommentRepository interface {
	Save(ctx context.Context, comment *entity.Comment) error
	Delete(ctx context.Context, comment *entity.Comment) error
}
