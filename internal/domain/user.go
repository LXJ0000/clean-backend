package domain

import (
	"context"
)

const (
	UserAuthorization = "x-user-id"
)

type User struct {
	Model
	UserID   int64  `json:"user_id" gorm:"primaryKey"`
	UserName string `json:"user_name" gorm:"unique"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type UserRepository interface {
	Create(c context.Context, u User) error
	GetByEmail(c context.Context, email string) (User, error)
	GetByID(c context.Context, id int64) (User, error)
	// FindOrCreate(c context.Context, phone string) (User, error)
}

type UserService interface {
	Signup(c context.Context, req RegisterReq) (Response, error)
	Login(c context.Context, req LoginReq) (Response, error)
	Detail(c context.Context, req UserDetailReq) (Response, error)
}

func (User) TableName() string {
	return `user`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterReq struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserDetailReq struct {
	UserID int64 `json:"user_id"`
}
