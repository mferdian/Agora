package controller

import (
	"Agora/constants"
	"Agora/dto"
	"Agora/logging"
	"Agora/service"
	"Agora/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	ICommentController interface {
		CreateComment(ctx *gin.Context)
		DeleteComment(ctx *gin.Context)
	}

	CommentController struct {
		commentService service.ICommentService
	}
)

func NewCommentController(commentService service.ICommentService) *CommentController {
	return &CommentController{
		commentService: commentService,
	}
}

func (cc *CommentController) CreateComment(ctx *gin.Context) {
	var req dto.CreateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_PROSES_REQUEST, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	resData, err := cc.commentService.CreateComment(ctx.Request.Context(), req)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_CREATE_COMMENT)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_CREATE_COMMENT, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_CREATE_COMMENT+": %s", resData.Content)
	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_CREATE_COMMENT, resData)
	ctx.JSON(http.StatusCreated, res)
}
