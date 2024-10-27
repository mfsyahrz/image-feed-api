package postgres

import (
	"database/sql"
	"fmt"

	"github.com/mfsyahrz/image_feed_api/internal/common/util"
	"github.com/mfsyahrz/image_feed_api/internal/domain/entity"
	"github.com/mfsyahrz/image_feed_api/internal/domain/repository"
)

type getInput struct {
	*repository.GetPostInput
}

func (i *getInput) queryArgs() ([]string, []interface{}) {
	var conds []string
	var args []interface{}

	conds = append(conds, "1 = 1")

	if len(i.IDs) > 0 {
		conds = append(conds, fmt.Sprintf("id in (%s)", util.JoinNumbers(i.IDs)))
	}

	if i.Cursor != nil {
		conds = append(conds, "(comment_count, created_date, id) < ($1, $2, $3)")
		args = append(args, i.Cursor.CommentCount, i.Cursor.CreatedAt, i.Cursor.PostID)
	}

	return conds, args
}

type postDTO struct {
	commentDTO
	ID            sql.NullInt64  `db:"id"`
	Caption       sql.NullString `db:"caption"`
	Creator       sql.NullString `db:"creator"`
	SrcImg        sql.NullString `db:"src_image"`
	DisplayImg    sql.NullString `db:"display_image"`
	CommentCount  sql.NullInt64  `db:"comment_count"`
	CreatedAt     sql.NullTime   `db:"created_date"`
}

func (d *postDTO) toPost() *entity.Post {
	return &entity.Post{
		ID:        d.ID.Int64,
		Caption:   d.Caption.String,
		Creator:   d.Creator.String,
		CreatedAt: d.CreatedAt.Time,
		DisplayImg: d.DisplayImg.String,
		SrcImg: d.SrcImg.String,
		CommentCount: d.CommentCount.Int64,
	}
}

type commentDTO struct {
	ID        sql.NullInt64  `db:"comment_id"`
	PostID    sql.NullInt64  `db:"post_id"`
	Content   sql.NullString `db:"comment_content"`
	Creator   sql.NullString `db:"comment_creator"`
	CreatedAt sql.NullTime   `db:"comment_created_date"`
}

func (d *commentDTO) toComment() *entity.Comment {
	return &entity.Comment{
		ID:        d.ID.Int64,
		PostID:    d.PostID.Int64,
		Creator:   d.Creator.String,
		Content:   d.Content.String,
		CreatedAt: d.CreatedAt.Time,
	}
}
