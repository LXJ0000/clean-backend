package handler

import (
	"net/http"

	"github.com/LXJ0000/clean-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService domain.UserService
}

func (u *UserHandler) Detail(c *gin.Context) {
	userID := c.MustGet(domain.UserAuthorization).(int64)
	resp, err := u.UserService.Detail(c, domain.UserDetailReq{UserID: userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (u *UserHandler) Login(c *gin.Context) {
	var req domain.LoginReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	resp, err := u.UserService.Login(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (u *UserHandler) Signup(c *gin.Context) {
	var req domain.RegisterReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	resp, err := u.UserService.Signup(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}
