package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/LXJ0000/clean-backend/bootstrap"
	"github.com/LXJ0000/clean-backend/internal/domain"
	snowflakeutil "github.com/LXJ0000/clean-backend/utils/snowflake"
)

type postService struct {
	postRepo           domain.PostRepository
	interactionService domain.InteractionService
	contextTimeout     time.Duration
	env                *bootstrap.Env
}

func NewPostService(postRepo domain.PostRepository, timeout time.Duration, env *bootstrap.Env, interactionService domain.InteractionService) domain.PostService {
	return &postService{
		postRepo:           postRepo,
		contextTimeout:     timeout,
		env:                env,
		interactionService: interactionService,
	}
}

func (u *postService) Create(c context.Context, req domain.PostCreateRequest) (domain.Response, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	post := domain.Post{
		PostID:   snowflakeutil.GenID(),
		Title:    req.Title,
		Abstract: req.Abstract,
		AuthorID: req.AuthorID,
		Content:  req.Content,
	}
	if err := u.postRepo.Create(ctx, post); err != nil {
		return domain.ErrorResp("Create post failed", err), err
	}
	return domain.SuccesResp(post), nil
}

func (u *postService) ReaderList(c context.Context, req domain.PostListRequest) (domain.Response, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	filter := map[string]interface{}{"status": domain.PostStatusPublish}
	if req.AuthorID != 0 {
		filter["author_id"] = req.AuthorID
	}

	items, err := u.postRepo.List(ctx, filter, req.Page, req.Size)
	if err != nil {
		return domain.ErrorResp("Get post list failed", err), err
	}
	cnt, err := u.postRepo.Count(ctx, filter)
	if err != nil {
		return domain.ErrorResp("Get post count failed", err), err
	}
	return domain.SuccesResp(map[string]interface{}{
		"cnt":       cnt,
		"post_list": items,
	}), nil
}

func (u *postService) WriterList(c context.Context, req domain.PostListRequest) (domain.Response, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if req.AuthorID == 0 {
		return domain.ErrorResp("AuthorID is required", nil), errors.New("AuthorID is required")
	}

	filter := map[string]interface{}{"author_id": req.AuthorID}
	items, err := u.postRepo.List(ctx, filter, req.Page, req.Size)
	if err != nil {
		return domain.ErrorResp("Get post list failed", err), err
	}
	cnt, err := u.postRepo.Count(ctx, filter)
	if err != nil {
		return domain.ErrorResp("Get post count failed", err), err
	}
	return domain.SuccesResp(map[string]interface{}{
		"cnt":       cnt,
		"post_list": items,
	}), nil
}

func (u *postService) Detail(c context.Context, req domain.PostDetailRequest) (domain.Response, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	post, err := u.postRepo.GetByID(ctx, req.PostID)
	if err != nil {
		return domain.ErrorResp("Get post detail failed", err), err
	}
	// go func() {
	// TODO kafka consume read count
	if err := u.interactionService.IncrReadCount(ctx, domain.BizPost, req.PostID); err != nil {
		slog.Warn("Incr read count failed", "Error", err)
	}
	// }()
	return domain.SuccesResp(map[string]interface{}{
		"post_detail": post,
		
	}), nil
}
