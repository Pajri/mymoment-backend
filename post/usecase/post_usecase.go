package usecase

import (
	"time"

	"github.com/pajri/personal-backend/domain"
)

type PostUsecase struct {
	postRepo    domain.IPostRepository
	accountRepo domain.IAccountRepository
}

func NewPostUseCase(postRepository domain.IPostRepository) *PostUsecase {
	return &PostUsecase{
		postRepo: postRepository,
	}
}

func (uc PostUsecase) InsertPost(post domain.Post) (*domain.Post, error) {
	post.Date = time.Now()
	err := uc.postRepo.InsertPost(post)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (uc PostUsecase) PostListing(accountID string, limit uint64, date time.Time) ([]domain.Post, error) {
	var postList []domain.Post

	var filter domain.PostFilter
	filter.AccountID = accountID
	filter.Limit = limit
	filter.Date = date
	postList, err := uc.postRepo.PostList(filter)
	return postList, err
}

func (uc PostUsecase) DeletePost(postID, accountID string) error {
	err := uc.postRepo.DeletePost(postID, accountID)
	return err
}
