package domain

import (
	"context"
)

const (
	KeyPostTopN = "post_topN"

	PostStatusHide uint8 = iota
	PostStatusPublish
)

type Post struct {
	Model
	PostID int64 `json:"post_id,string" gorm:"primaryKey"`

	Title    string
	Abstract string
	Content  string
	AuthorID int64
	Status   uint8 `gorm:"default:1"`
}

func (Post) TableName() string {
	return `post`
}

//go:generate mockgen -source=./post.go -destination=./mock/post.go -package=domain_mock
type PostRepository interface {
	Create(c context.Context, post Post) error
	GetByID(c context.Context, id int64) (Post, error)
	List(c context.Context, filter interface{}, page, size int) ([]Post, error)
	Count(c context.Context, filter interface{}) (int64, error)
}

type PostService interface {
	Create(c context.Context, req PostCreateRequest) (Response, error)
	ReaderList(c context.Context, req PostListRequest) (Response, error)
	WriterList(c context.Context, req PostListRequest) (Response, error)
	Detail(c context.Context, req PostDetailRequest) (Response, error)
	//ReplaceTopN(c context.Context, items []Post, expiration time.Duration) error
}

type PostListRequest struct {
	Page     int   `json:"page" form:"page"`
	Size     int   `json:"size" form:"size"`
	AuthorID int64 `json:"author_id" form:"author_id"`
}

type PostCreateRequest struct {
	Title    string `json:"title" form:"title" binding:"required"`
	Abstract string `json:"abstract" form:"abstract" binding:"required"`
	Content  string `json:"content" form:"content" binding:"required"`
	AuthorID int64  `json:"author_id" form:"author_id" binding:"required"`
}

type PostDetailRequest struct {
	PostID int64 `json:"post_id" form:"post_id"`
}
