package controller

import (
	"Agora/constants"
	"Agora/dto"
	"Agora/logging"
	"Agora/service"
	"Agora/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	IProposalController interface {
		CreateProposal(ctx *gin.Context)
		GetAllProposal(ctx *gin.Context)
		GetProposalById(ctx *gin.Context)
		UpdateProposal(ctx *gin.Context)
		DeleteProposal(ctx *gin.Context)
	}

	ProposalController struct {
		proposalService service.IProposalService
	}
)

func NewProposalController(proposalService service.IProposalService) *ProposalController {
	return &ProposalController{
		proposalService: proposalService,
	}
}

func (ps *ProposalController) CreateProposal(ctx *gin.Context) {
	var req dto.CreateProposalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_PROSES_REQUEST, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	resData, err := ps.proposalService.CreateProposal(ctx.Request.Context(), req)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_CREATE_PROPOSAL)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_CREATE_PROPOSAL, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_CREATE_PROPOSAL+": %s", resData.Title)
	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_CREATE_PROPOSAL, resData)
	ctx.JSON(http.StatusCreated, res)
}
func (ps *ProposalController) GetAllProposal(ctx *gin.Context) {
	paginationParam := ctx.DefaultQuery("pagination", "true")
	usePagination := paginationParam != "false"

	if !usePagination {
		result, err := ps.proposalService.GetAllProposal(ctx.Request.Context())
		if err != nil {
			logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_GET_LIST_PROPOSAL)
			res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_LIST_PROPOSAL, err.Error(), nil)
			ctx.JSON(http.StatusBadRequest, res)
			return
		}
		logging.Log.Info(constants.MESSAGE_SUCCESS_GET_LIST_PROPOSAL + " (no pagination)")
		res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_GET_LIST_PROPOSAL, result)
		ctx.JSON(http.StatusOK, res)
		return
	}

	var query dto.ProposalPaginationRequest
	if err := ctx.ShouldBindQuery(&query); err != nil {
		logging.Log.WithError(err).Warn(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := ps.proposalService.GetAllProposalWithPagination(ctx.Request.Context(), query)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_GET_LIST_PROPOSAL)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_LIST_PROPOSAL, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_GET_LIST_PROPOSAL+": page %d", query.Page)
	res := utils.Response{
		Status:   true,
		Messsage: constants.MESSAGE_SUCCESS_GET_LIST_PROPOSAL,
		Data:     result.Data,
		Meta:     result.PaginationResponse,
	}
	ctx.JSON(http.StatusOK, res)
}
func (ps *ProposalController) GetProposalById(ctx *gin.Context) {
	idStr := ctx.Param("id")

	if _, err := uuid.Parse(idStr); err != nil {
		logging.Log.WithError(err).Warn(constants.MESSAGE_FAILED_UUID_FORMAT)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UUID_FORMAT, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ps.proposalService.GetProposalByID(ctx.Request.Context(), idStr)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_GET_DETAIL_PROPOSAL)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DETAIL_PROPOSAL, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusNotFound, res)
		return
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_GET_DETAIL_PROPOSAL+": %s", idStr)
	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_GET_DETAIL_PROPOSAL, result)
	ctx.JSON(http.StatusOK, res)
}
func (ps *ProposalController) UpdateProposal(ctx *gin.Context) {
	idParam := ctx.Param("id")
	if _, err := uuid.Parse(idParam); err != nil {
		logging.Log.WithError(err).Warn(constants.MESSAGE_FAILED_UUID_FORMAT)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UUID_FORMAT, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	var payload dto.UpdateProposalRequest
	payload.ID = idParam

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logging.Log.WithError(err).Warn(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := ps.proposalService.UpdateProposal(ctx.Request.Context(), payload)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_UPDATE_PROPOSAL)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UPDATE_PROPOSAL, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_UPDATE_PROPOSAL+": %s", result.ID)
	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_UPDATE_PROPOSAL, result)
	ctx.JSON(http.StatusOK, res)
}
func (ps *ProposalController) DeleteProposal(ctx *gin.Context) {
	idParam := ctx.Param("id")
	if _, err := uuid.Parse(idParam); err != nil {
		logging.Log.WithError(err).Warn(constants.MESSAGE_FAILED_UUID_FORMAT)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UUID_FORMAT, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	payload := dto.DeleteProposalRequest{ID: idParam}

	result, err := ps.proposalService.DeleteProposal(ctx.Request.Context(), payload)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_DELETE_PROPOSAL)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_DELETE_PROPOSAL, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_DELETE_PROPOSAL+": %s", idParam)
	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_DELETE_PROPOSAL, result)
	ctx.JSON(http.StatusOK, res)
}
