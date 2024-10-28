package entity

import (
	"encoding/json"
	"time"
)

const (
	MaxComments = 2
	MaxPosts    = 10
)

type Comments []*Comment

func (c Comments) Len() int {
	return len(c)
}

type Comment struct {
	ID        int64     `json:"id" db:"comment_id"`
	PostID    int64     `json:"postID" db:"post_id"`
	Content   string    `json:"content" db:"comment_content"`
	Creator   string    `json:"creator" db:"comment_creator"`
	CreatedAt time.Time `json:"createdAt" db:"comment_created_date"`
}

func (t Comment) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Comment) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &t)
}
