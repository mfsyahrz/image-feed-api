package comment

import (
	"context"
	"fmt"

	"github.com/mfsyahrz/image_feed_api/internal/common/logger"
	"github.com/mfsyahrz/image_feed_api/internal/domain/entity"
	"github.com/mfsyahrz/image_feed_api/internal/domain/repository"
)

type CommentService interface {
	Save(ctx context.Context, comment *entity.Comment) error
	Delete(ctx context.Context, comment *entity.Comment) error
}

type commentService struct {
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
}

func NewCommentService(commentRepo repository.CommentRepository, postRepo repository.PostRepository) CommentService {
	return &commentService{commentRepo, postRepo}
}

func (s *commentService) Save(ctx context.Context, comment *entity.Comment) error {
	_, err := s.postRepo.FetchOne(ctx, repository.GetPostInput{
		IDs: []int64{comment.PostID},
	})
	if err != nil {
		logger.FromCtx(ctx).Error(err.Error())
		return fmt.Errorf("unable to get postID %d. error: %s", comment.PostID, err.Error())
	}

	if err := s.commentRepo.Save(ctx, comment); err != nil {
		logger.FromCtx(ctx).Error(err.Error())
		return fmt.Errorf("unable to save comment. error: %s", err.Error())
	}
	return nil
}

func (s *commentService) Delete(ctx context.Context, comment *entity.Comment) error {
	_, err := s.postRepo.FetchOne(ctx, repository.GetPostInput{
		IDs: []int64{comment.PostID},
	})
	if err != nil {
		logger.FromCtx(ctx).Error(err.Error())
		return fmt.Errorf("unable to get postID %d. error: %s", comment.PostID, err.Error())
	}

	if err := s.commentRepo.Delete(ctx, comment); err != nil {
		logger.FromCtx(ctx).Error(err.Error())
		return fmt.Errorf("unable to delete comment. error: %+v", err.Error())
	}

	return nil
}
