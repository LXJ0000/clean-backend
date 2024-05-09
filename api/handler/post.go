package handler

import (
	"net/http"

	"github.com/LXJ0000/clean-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	PostService domain.PostService
}

func (h *PostHandler) Publish(c *gin.Context) {

}

func (h *PostHandler) Create(c *gin.Context) {
	var q domain.PostCreateRequest
	if err := c.ShouldBind(&q); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.PostService.Create(c, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *PostHandler) ReaderList(c *gin.Context) {
	var q domain.PostListRequest
	if err := c.ShouldBind(&q); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.PostService.ReaderList(c, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *PostHandler) WriterList(c *gin.Context) {
	var q domain.PostListRequest
	if err := c.ShouldBind(&q); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.PostService.WriterList(c, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *PostHandler) Detail(c *gin.Context) {
	var q domain.PostDetailRequest
	if err := c.ShouldBind(&q); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.PostService.Detail(c, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}
