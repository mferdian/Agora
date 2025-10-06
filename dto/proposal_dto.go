package dto

import (
	"Agora/model"

	"github.com/google/uuid"
)

type (

	// Response
	ProposalResponse struct {
		ID          uuid.UUID `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
	}

	UserCompactResponse struct {
		ID    uuid.UUID `json:"id"`
		Name  string    `json:"name"`
		Email string    `json:"email"`
	}

	CreateProposalResponse struct {
		ID          uuid.UUID           `json:"id"`
		Title       string              `json:"title"`
		Description string              `json:"description"`
		User        UserCompactResponse `json:"user"`
	}

	// Request
	CreateProposalRequest struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
	}
	UpdateProposalRequest struct {
		ID          string  `json:"id"`
		Title       *string `json:"title,omitempty"`
		Description *string `json:"description,omitempty"`
	}

	DeleteProposalRequest struct {
		ID string `json:"-"`
	}

	// Pagination
	ProposalPaginationRequest struct {
		PaginationRequest
		ID string `form:"id"`
	}

	ProposalPaginationResponse struct {
		PaginationResponse
		Data []ProposalResponse `json:"data"`
	}

	ProposalPaginationRepositoryResponse struct {
		PaginationResponse
		Proposals []model.Proposal
	}
)
