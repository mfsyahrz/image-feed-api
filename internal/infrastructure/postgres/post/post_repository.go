package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/mfsyahrz/image_feed_api/internal/common/logger"
	"github.com/mfsyahrz/image_feed_api/internal/common/pagination"
	"github.com/mfsyahrz/image_feed_api/internal/common/util"
	"github.com/mfsyahrz/image_feed_api/internal/domain/entity"
	"github.com/mfsyahrz/image_feed_api/internal/domain/repository"
)

const (
	findPostsQuery = `SELECT p.id AS id, p.caption, p.creator, p.display_image, p.comment_count, p.created_date,
					  c.comment_id, c.post_id, comment_creator, comment_content, comment_created_date
					  FROM (
							SELECT id, caption, creator, display_image, comment_count, created_date
							FROM posts
							%[1]s
							ORDER BY comment_count DESC, created_date DESC, id ASC
							LIMIT %[2]d
					   ) AS p
					   LEFT JOIN (
							SELECT id AS comment_id, post_id, creator as comment_creator, content as comment_content, created_date AS comment_created_date,
										ROW_NUMBER() OVER (PARTITION BY post_id ORDER BY created_date DESC) AS rn
							FROM comments
					   ) AS c ON p.id = c.post_id AND c.rn <= %[3]d
					  ORDER BY p.comment_count DESC, p.created_date DESC, p.id ASC, comment_created_date DESC;
				`

	insertPostQuery = `insert into posts (caption, creator, src_image, display_image) 
					   values ($1, $2, $3, $4) RETURNING id`
)

type postRepo struct {
	db *sqlx.DB
}

func NewPostRepo(db *sqlx.DB) repository.PostRepository {
	return &postRepo{db}
}

func (r *postRepo) FetchPaginated(ctx context.Context, input repository.GetPostInput) (entity.Posts, *pagination.PostCursor, error) {
	posts, err := r.fetch(ctx, input)
	if err != nil {
		return nil, nil, err
	}

	if posts.Len() < input.Limit {
		return posts, nil, nil
	}

	return posts, pagination.FromPosts(posts), nil
}

func (r *postRepo) FetchOne(ctx context.Context, input repository.GetPostInput) (*entity.Post, error) {
	input.Limit = 1

	posts, err := r.fetch(ctx, input)
	if err != nil {
		return nil, err
	}

	if posts.Len() != 1 {
		return nil, fmt.Errorf("no post found for query")
	}

	return posts.GetFirst(), nil
}

func (r *postRepo) fetch(ctx context.Context, input repository.GetPostInput) (entity.Posts, error) {
	logger.FromCtx(ctx).Info("Fetching posts for input: %s", util.PrettyPrint(input))

	req := getInput{&input}
	conds, args := req.queryArgs()
	limit := util.DefaultIfZero(input.Limit, entity.MaxPosts)
	commentLimit := util.DefaultIfZero(input.CommentLimit, entity.MaxComments)
	query := fmt.Sprintf(findPostsQuery, "WHERE "+strings.Join(conds, " AND "), limit, commentLimit)

	var results []*postDTO

	err := r.db.Select(&results, query, args...)
	if err != nil {
		logger.FromCtx(ctx).Error("failed to fetch posts %s", err.Error())
		return nil, err
	}

	return r.decodePost(results)
}

func (r *postRepo) decodePost(rows []*postDTO) (entity.Posts, error) {
	postMap := make(map[int64]*entity.Post)
	var posts entity.Posts
	for _, row := range rows {
		postID := row.ID.Int64

		if _, exists := postMap[postID]; !exists {
			post := row.toPost()
			post.Comments = []*entity.Comment{}
			posts = append(posts, post)
			postMap[postID] = post
		}

		if row.commentDTO.ID.Int64 != 0 {
			comment := row.toComment()
			postMap[postID].Comments = append(postMap[postID].Comments, comment)
		}
	}

	return posts, nil
}

func (r *postRepo) Save(ctx context.Context, post *entity.Post) error {
	if post.CommentLen() > entity.MaxComments {
		return fmt.Errorf("post comments exceed comment limit of %d", entity.MaxComments)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			logger.FromCtx(ctx).Error("failed to save post %s", err.Error())
		}
	}()

	err = tx.QueryRow(insertPostQuery, post.Caption, post.Creator, post.SrcImg, post.DisplayImg).Scan(&post.ID)
	if err != nil {
		return err
	}

	for _, comment := range post.Comments {
		err := tx.QueryRow(insertPostQuery, post.Caption, post.Creator, post.SrcImg, post.DisplayImg).Scan(&comment.ID)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
