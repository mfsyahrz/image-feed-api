package pagination

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mfsyahrz/image_feed_api/internal/domain/entity"
)

// PostCursor defines the cursor for paginating posts
type PostCursor struct {
	CommentCount int64
	PostID       int64
	CreatedAt    time.Time
}

// Encode encodes the PostCursor into a string
func (c *PostCursor) Encode() string {
	if c == nil {
		return ""
	}
	return fmt.Sprintf("%d-%d-%d", c.CommentCount, c.PostID, c.CreatedAt.Unix())
}

// DecodePostCursor decodes a cursor string into a PostCursor
func DecodePostCursor(cursorStr string) (*PostCursor, error) {
	if cursorStr == "" {
		return nil, nil
	}

	parts := strings.Split(cursorStr, "-")
	if len(parts) != 3 {
		return nil, errors.New("invalid cursor format")
	}

	commentCount, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor comment count: %w", err)
	}

	postID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor postID: %w", err)
	}

	createdAtUnix, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid created_at timestamp: %w", err)
	}
	createdAt := time.Unix(createdAtUnix, 0)

	return &PostCursor{
		CommentCount: commentCount,
		PostID:       postID,
		CreatedAt:    createdAt,
	}, nil
}

// FromPosts converts a PostCursor from Posts domain
func FromPosts(posts entity.Posts) *PostCursor {
	lastPost := posts.GetLast()
	if lastPost == nil {
		return nil
	}

	return &PostCursor{
		CommentCount: lastPost.CommentCount,
		PostID:       lastPost.ID,
		CreatedAt:    lastPost.CreatedAt,
	}
}
