package repository

import (
	"Agora/dto"
	"Agora/model"
	"context"
	"math"
	"strings"

	"gorm.io/gorm"
)

type (
	IProposalRepository interface {
		CreateProposal(ctx context.Context, tx *gorm.DB, proposal model.Proposal) error
		GetAllProposalWithPagination(ctx context.Context, tx *gorm.DB, req dto.ProposalPaginationRequest) (dto.ProposalPaginationRepositoryResponse, error)
		GetAllProposal(ctx context.Context, tx *gorm.DB) ([]model.Proposal, error)
		GetProposalByID(ctx context.Context, tx *gorm.DB, id string) (model.Proposal, bool, error)
		UpdateProposal(ctx context.Context, tx *gorm.DB, proposal model.Proposal) error
		DeleteProposal(ctx context.Context, tx *gorm.DB, id string) error
	}

	ProposalRepository struct {
		db *gorm.DB
	}
)

func NewProposalRepository(db *gorm.DB) *ProposalRepository {
	return &ProposalRepository{
		db: db,
	}
}

func (pr *ProposalRepository) CreateProposal(ctx context.Context, tx *gorm.DB, proposal model.Proposal) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Create(&proposal).Error
}

func (pr *ProposalRepository) GetAllProposalWithPagination(ctx context.Context, tx *gorm.DB, req dto.ProposalPaginationRequest) (dto.ProposalPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = pr.db
	}

	var proposals []model.Proposal
	var err error
	var count int64

	if req.PaginationRequest.PerPage == 0 {
		req.PaginationRequest.PerPage = 10
	}

	if req.PaginationRequest.Page == 0 {
		req.PaginationRequest.Page = 1
	}

	query := tx.WithContext(ctx).Model(&model.Proposal{})

	if req.PaginationRequest.Search != "" {
		searchValue := "%" + strings.ToLower(req.PaginationRequest.Search) + "%"
		query = query.Where("LOWER(title) LIKE ?", searchValue)
	}

	if req.ID != "" {
		query = query.Where("id = ?", req.ID)
	}

	query = query.Preload("User")

	if err := query.Count(&count).Error; err != nil {
		return dto.ProposalPaginationRepositoryResponse{}, err
	}

	if err := query.Order("created_at DESC").Scopes(Paginate(req.PaginationRequest.Page, req.PaginationRequest.PerPage)).Find(&proposals).Error; err != nil {
		return dto.ProposalPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PaginationRequest.PerPage)))

	return dto.ProposalPaginationRepositoryResponse{
		Proposals: proposals,
		PaginationResponse: dto.PaginationResponse{
			Page:    req.PaginationRequest.Page,
			PerPage: req.PaginationRequest.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, err
}

func (pr *ProposalRepository) GetProposalByID(ctx context.Context, tx *gorm.DB, id string) (model.Proposal, bool, error) {
	if tx == nil {
		tx = pr.db
	}

	var proposal model.Proposal
	if err := tx.WithContext(ctx).Where("id = ?", id).Take(&proposal).Error; err != nil {
		return model.Proposal{}, false, err
	}

	return proposal, true, nil
}

func (pr *ProposalRepository) UpdateProposal(ctx context.Context, tx *gorm.DB, proposal model.Proposal) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Where("id = ?", proposal.ID).Updates(&proposal).Error
}

func (pr *ProposalRepository) GetAllProposal(ctx context.Context, tx *gorm.DB) ([]model.Proposal, error) {
	var proposals []model.Proposal

	err := pr.db.WithContext(ctx).Model(&model.Proposal{}).Preload("User").Find(&proposals).Error

	if err != nil {
		return nil, err
	}

	return proposals, nil
}

func (pr *ProposalRepository) DeleteProposal(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Where("id = ?", id).Delete(&model.Proposal{}).Error
}
