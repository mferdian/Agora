package service

import (
	"Agora/constants"
	"Agora/dto"
	"Agora/helpers"
	"Agora/logging"
	"Agora/model"
	"Agora/repository"
	"context"

	"github.com/google/uuid"
)

type (
	IProposalService interface {
		CreateProposal(ctx context.Context, req dto.CreateProposalRequest) (dto.CreateProposalResponse, error)
		GetAllProposalWithPagination(ctx context.Context, req dto.ProposalPaginationRequest) (dto.ProposalPaginationResponse, error)
		GetAllProposal(ctx context.Context) ([]dto.ProposalResponse, error)
		GetProposalByID(ctx context.Context, id string) (dto.ProposalResponse, error)
		UpdateProposal(ctx context.Context, req dto.UpdateProposalRequest) (dto.ProposalResponse, error)
		DeleteProposal(ctx context.Context, req dto.DeleteProposalRequest) (dto.ProposalResponse, error)
	}

	ProposalService struct {
		proposalRepo repository.IProposalRepository
		jwtService   InterfaceJWTService
	}
)

func NewProposalService(proposalRepo repository.IProposalRepository, jwtService InterfaceJWTService) *ProposalService {
	return &ProposalService{
		proposalRepo: proposalRepo,
		jwtService:   jwtService,
	}
}

func (ps *ProposalService) CreateProposal(ctx context.Context, req dto.CreateProposalRequest) (dto.CreateProposalResponse, error) {
	userID := helpers.GetUserID(ctx)
	if userID == "" {
		return dto.CreateProposalResponse{}, constants.ErrGetIDFromToken
	}

	if len(req.Title) < 2 {
		logging.Log.Warn(constants.MESSAGE_FAILED_CREATE_PROPOSAL + ": name too short")
		return dto.CreateProposalResponse{}, constants.ErrInvalidProposalName
	}

	proposal := model.Proposal{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		UserID:      uuid.MustParse(userID),
	}

	if err := ps.proposalRepo.CreateProposal(ctx, nil, proposal); err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_CREATE_PROPOSAL)
		return dto.CreateProposalResponse{}, constants.ErrCreateProposal
	}

	return dto.CreateProposalResponse{
		ID:          proposal.ID,
		Title:       proposal.Title,
		Description: proposal.Description,
		User: dto.UserCompactResponse{
			ID:    uuid.MustParse(userID),
			Name:  proposal.User.Name,
			Email: proposal.User.Email,
		},
	}, nil

}

func (ps *ProposalService) GetAllProposalWithPagination(ctx context.Context, req dto.ProposalPaginationRequest) (dto.ProposalPaginationResponse, error) {
	dataWithPaginate, err := ps.proposalRepo.GetAllProposalWithPagination(ctx, nil, req)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_GET_LIST_PROPOSAL)
		return dto.ProposalPaginationResponse{}, constants.ErrGetAllProposalWithPagination
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_GET_LIST_PROPOSAL+": page %d", req.Page)

	var datas []dto.ProposalResponse
	for _, proposal := range dataWithPaginate.Proposals {
		datas = append(datas, dto.ProposalResponse{
			ID:          proposal.ID,
			Title:       proposal.Title,
			Description: proposal.Description,
		})
	}

	return dto.ProposalPaginationResponse{
		Data: datas,
		PaginationResponse: dto.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}
func (ps *ProposalService) GetAllProposal(ctx context.Context) ([]dto.ProposalResponse, error) {
	proposals, err := ps.proposalRepo.GetAllProposal(ctx, nil)
	if err != nil {
		return nil, constants.ErrGetAllProposal
	}

	var datas []dto.ProposalResponse
	for _, proposal := range proposals {
		datas = append(datas, dto.ProposalResponse{
			ID:          proposal.ID,
			Title:       proposal.Title,
			Description: proposal.Description,
		})
	}

	return datas, nil
}
func (ps *ProposalService) GetProposalByID(ctx context.Context, id string) (dto.ProposalResponse, error) {
	if _, err := uuid.Parse(id); err != nil {
		logging.Log.Warn(constants.MESSAGE_FAILED_GET_DETAIL_PROPOSAL + ": invalid UUID")
		return dto.ProposalResponse{}, constants.ErrInvalidUUID
	}

	proposal, _, err := ps.proposalRepo.GetProposalByID(ctx, nil, id)
	if err != nil {
		logging.Log.WithError(err).WithField("id", id).Error(constants.MESSAGE_FAILED_GET_DETAIL_USER)
		return dto.ProposalResponse{}, constants.ErrGetProposalByID
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_GET_DETAIL_PROPOSAL+": %s", id)

	return dto.ProposalResponse{
		ID:          proposal.ID,
		Title:       proposal.Title,
		Description: proposal.Description,
	}, nil
}
func (ps *ProposalService) UpdateProposal(ctx context.Context, req dto.UpdateProposalRequest) (dto.ProposalResponse, error) {
	proposal, _, err := ps.proposalRepo.GetProposalByID(ctx, nil, req.ID)
	if err != nil {
		logging.Log.WithError(err).WithField("id", req.ID).Error(constants.MESSAGE_FAILED_UPDATE_PROPOSAL)
		return dto.ProposalResponse{}, constants.ErrGetProposalByID
	}

	if req.Title != nil && len(*req.Title) < 2 {
		logging.Log.Warn(constants.MESSAGE_FAILED_UPDATE_PROPOSAL + ": invalid title")
		return dto.ProposalResponse{}, constants.ErrInvalidName
	} else if req.Title != nil {
		proposal.Title = *req.Title
	}
	if req.Description != nil && len(*req.Description) < 2 {
		logging.Log.Warn(constants.MESSAGE_FAILED_UPDATE_PROPOSAL + ": invalid description")
		return dto.ProposalResponse{}, constants.ErrInvalidName
	} else if req.Description != nil {
		proposal.Description = *req.Description
	}

	err = ps.proposalRepo.UpdateProposal(ctx, nil, proposal)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_UPDATE_PROPOSAL)
		return dto.ProposalResponse{}, constants.ErrUpdateProposal
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_UPDATE_PROPOSAL+": %s", proposal.ID)

	return dto.ProposalResponse{
		ID:          proposal.ID,
		Title:       proposal.Title,
		Description: proposal.Description,
	}, nil
}

func (ps *ProposalService) DeleteProposal(ctx context.Context, req dto.DeleteProposalRequest) (dto.ProposalResponse, error) {
	proposal, _, err := ps.proposalRepo.GetProposalByID(ctx, nil, req.ID)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_DELETE_PROPOSAL)
		return dto.ProposalResponse{}, constants.ErrGetProposalByID
	}

	err = ps.proposalRepo.DeleteProposal(ctx, nil, req.ID)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_DELETE_USER)
		return dto.ProposalResponse{}, constants.ErrDeleteProposal
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_DELETE_USER+": %s", req.ID)

	return dto.ProposalResponse{
		ID:          proposal.ID,
		Title:       proposal.Title,
		Description: proposal.Description,
	}, nil
}
