package service

import (
	"context"
	"errors"
	"time"

	"github.com/LXJ0000/clean-backend/bootstrap"
	"github.com/LXJ0000/clean-backend/internal/domain"
	snowflakeutil "github.com/LXJ0000/clean-backend/utils/snowflake"
	tokenutil "github.com/LXJ0000/clean-backend/utils/token"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo       domain.UserRepository
	Env            *bootstrap.Env
	contextTimeout time.Duration
}

func NewUserService(userRepo domain.UserRepository, Env *bootstrap.Env, contextTimeout time.Duration) domain.UserService {
	return &userService{
		userRepo:       userRepo,
		Env:            Env,
		contextTimeout: contextTimeout,
	}
}

func (u *userService) Detail(c context.Context, req domain.UserDetailReq) (domain.Response, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return domain.ErrorResp("User Not Found With DB Error", err), err
	}
	return domain.SuccesResp(user), nil
}

func (u *userService) Login(c context.Context, req domain.LoginReq) (domain.Response, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return domain.ErrorResp("User not found with the given email", err), err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return domain.ErrorResp("Invalid credentials", nil), errors.New("invalid credentials")
	}

	accessToken, err := u.createAccessToken(user, u.Env.AccessTokenSecret, u.Env.AccessTokenExpiryHour)
	if err != nil {
		return domain.ErrorResp(err.Error(), err), err
	}

	refreshToken, err := u.createRefreshToken(user, u.Env.RefreshTokenSecret, u.Env.RefreshTokenExpiryHour)
	if err != nil {
		return domain.ErrorResp(err.Error(), err), err
	}

	resp := domain.Response{
		Code:    0,
		Message: "success",
		Data: map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"user_detail":   user,
		},
	}
	return resp, nil
}

func (u *userService) Signup(c context.Context, req domain.RegisterReq) (domain.Response, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if _, err := u.userRepo.GetByEmail(c, req.Email); err == nil {
		return domain.ErrorResp("User already exists with the given email", nil), errors.New("user already exists with the given email")
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return domain.ErrorResp("User already exists with the given email", err), err
	}

	req.Password = string(encryptedPassword)

	user := domain.User{
		UserID:   snowflakeutil.GenID(),
		UserName: req.UserName,
		Email:    req.Email,
		Password: req.Password,
	}

	err = u.userRepo.Create(ctx, user)
	if err != nil {
		return domain.ErrorResp("Create User Fail With DB Error", err), err
	}

	accessToken, err := u.createAccessToken(user, u.Env.AccessTokenSecret, u.Env.AccessTokenExpiryHour)
	if err != nil {
		return domain.ErrorResp("Create AccessToken Fail", err), err
	}

	refreshToken, err := u.createRefreshToken(user, u.Env.RefreshTokenSecret, u.Env.RefreshTokenExpiryHour)
	if err != nil {
		return domain.ErrorResp("Create RefreshToken Fail", err), err
	}

	return domain.SuccesResp(map[string]interface{}{
		"access_tone":   accessToken,
		"refresh_token": refreshToken,
		"user_detail":   user,
	}), nil
}

func (u *userService) createAccessToken(user domain.User, secret string, expiry int) (accessToken string, err error) {
	return tokenutil.CreateAccessToken(user, secret, expiry)
}

func (u *userService) createRefreshToken(user domain.User, secret string, expiry int) (refreshToken string, err error) {
	return tokenutil.CreateRefreshToken(user, secret, expiry)
}
