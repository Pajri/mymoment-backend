package domain

import (
	"time"
)

type Post struct {
	PostID      string    `json:"post_id"`
	Content     string    `json:"content"`
	ImageURL    string    `json:"image_url"`
	Date        time.Time `json:"date"`
	LastUpdated time.Time `json:"last_updated"`
	AccountID   string    `json:"account_id"`
	Account     Account   `json:"-"`
}

type IPostRepository interface {
	InsertPost(post Post) error
	DeletePost(postID, accountID string) error
	PostList(filter PostFilter) ([]Post, error)
	GetPost(filter PostFilter) (*Post, error)
}

type IPostUsecase interface {
	InsertPost(post Post) (*Post, error)
	DeletePost(postID, accountID string) error
	PostListing(accountId string, limit uint64, date time.Time) ([]Post, error)
}

type PostFilter struct {
	PostID    string
	AccountID string
	Date      time.Time
	Limit     uint64
}
