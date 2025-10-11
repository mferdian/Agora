package controller

import (
	"Agora/helpers"
	"Agora/model"
	"Agora/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	ICommentController interface {
		CreateComment(ctx *gin.Context)
	}

	CommentController struct {
		commentService service.ICommentService
	}
)

func NewCommentController(s service.ICommentService) *CommentController {
	return &CommentController{commentService: s}
}

func (c *CommentController) CreateComment(ctx *gin.Context) {
	var req model.Comment
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := helpers.GetUserID(ctx)
	comment, err := c.commentService.CreateComment(userID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, comment)
}
