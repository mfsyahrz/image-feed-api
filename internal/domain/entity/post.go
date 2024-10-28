package entity

import (
	"encoding/json"
	"time"
)

type Posts []*Post

func (p Posts) GetIDs() []int64 {
	ids := make([]int64, len(p))

	for i, data := range p {
		ids[i] = data.ID
	}

	return ids
}

func (p Posts) Len() int {
	return len(p)
}

func (p Posts) GetLast() *Post {
	if len(p) == 0 {
		return nil
	}

	return p[p.Len()-1]
}

func (p Posts) GetFirst() *Post {
	if len(p) == 0 {
		return nil
	}

	return p[0]
}

type Post struct {
	ID            int64     `json:"id" db:"id"`
	Caption       string    `json:"caption"  db:"caption"`
	Creator       string    `json:"creator"  db:"creator"`
	SrcImg        string    `json:"-"  db:"src_image"`
	DisplayImg    string    `json:"-" db:"display_image"`
	DisplayImgURL string    `json:"displayImgURL" db:"-"`
	Comments      Comments  `json:"comments,omitempty"  db:"-"`
	CommentCount  int64     `json:"commentCount" db:"comment_count"`
	CreatedAt     time.Time `json:"createdAt"  db:"created_date"`
}

func (t *Post) CommentLen() int {
	return t.Comments.Len()
}

func (t *Post) SetDisplayImgURL(baseURL string) {
	t.DisplayImgURL = baseURL + t.DisplayImg
}

func (t Post) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Post) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}
