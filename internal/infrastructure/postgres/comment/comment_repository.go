package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/mfsyahrz/image_feed_api/internal/common/logger"
	"github.com/mfsyahrz/image_feed_api/internal/common/util"
	"github.com/mfsyahrz/image_feed_api/internal/domain/entity"
	"github.com/mfsyahrz/image_feed_api/internal/domain/repository"
)

const (
	insertCommentQuery = `insert into comments (post_id, content, creator) 
		values ($1, $2, $3) RETURNING id`
	deleteCommentQuery = `DELETE from comments WHERE id = $1`
	fetchCommentsQuery = `SELECT * FROM comments where post_id in (%s)`
	incrCommentCounterQuery = `UPDATE posts SET comment_count = comment_count + 1 WHERE id = $1`
	decrCommentCounterQuery = `UPDATE posts SET comment_count = comment_count - 1 WHERE id = $1`
)

type commentRepo struct {
	db *sqlx.DB
}

func NewCommentRepo(db *sqlx.DB) repository.CommentRepository {
	return &commentRepo{db}
}

func (r *commentRepo) Fetch(ctx context.Context, input repository.GetCommentInput) (entity.Comments, error) {
	logger.FromCtx(ctx).Info("Fetching comments with input %s", util.PrettyPrint(input))

	var comments entity.Comments

	query := fmt.Sprintf(fetchCommentsQuery, util.JoinNumbers(input.PostIDs))
	err := r.db.Select(&comments, query)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *commentRepo) Save(ctx context.Context, comment *entity.Comment) error {
	logger.FromCtx(ctx).Info("saving comment for postID %d", comment.PostID)

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			logger.FromCtx(ctx).Error("failed to save comment %s", err.Error())
		}
	}()

	err = tx.QueryRow(insertCommentQuery, comment.PostID, comment.Content, comment.Creator).Scan(&comment.ID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(incrCommentCounterQuery, comment.PostID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *commentRepo) Delete(ctx context.Context, comment *entity.Comment) error {
	logger.FromCtx(ctx).Info("deleting comment %d for postID %d", comment.ID, comment.PostID)

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			logger.FromCtx(ctx).Error("failed to delete comment %s", err.Error())
		}
	}()

	res, err := tx.Exec(deleteCommentQuery, comment.ID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		err = errors.New("no rows affected")
		return err
	}

	_, err = tx.Exec(decrCommentCounterQuery, comment.PostID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
