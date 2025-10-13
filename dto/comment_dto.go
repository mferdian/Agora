package dto

import "github.com/google/uuid"

type (
	// Response
	ProposalCommentResponse struct {
		ID    uuid.UUID `json:"id"`
		Title string    `json:"title"`
	}
	UserCommentResponse struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}
	CommentResponse struct {
		ID       uuid.UUID               `json:"id"`
		Proposal ProposalCommentResponse `json:"proposal"`
		User     UserCommentResponse    `json:"user"`
		Content  string                  `json:"comment_content"`
	}

	// request
	CreateCommentRequest struct {
		ProposalID uuid.UUID `json:"proposal_id"`
		Content   string    `json:"content"`
	}
	DeleteCommentRequest struct {
		CommentID string `json:"-"`
	}
)
