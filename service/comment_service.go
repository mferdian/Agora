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
	ICommentService interface {
		CreateComment(ctx context.Context, req dto.CreateCommentRequest) (dto.CommentResponse, error)
		DeleteComment(ctx context.Context, req dto.DeleteCommentRequest) (dto.CommentResponse, error)
	}

	CommentService struct {
		commentRepo  repository.ICommentRepository
		userRepo     repository.IUserRepository
		proposalRepo repository.IProposalRepository
		jwtService   IJWTService
	}
)

func NewCommentService(commentRepo repository.ICommentRepository, userRepo repository.IUserRepository, proposalRepo repository.IProposalRepository, jwtService IJWTService) *CommentService {
	return &CommentService{
		commentRepo:  commentRepo,
		userRepo:     userRepo,
		proposalRepo: proposalRepo,
		jwtService:   jwtService,
	}
}

func (cs *CommentService) CreateComment(ctx context.Context, req dto.CreateCommentRequest) (dto.CommentResponse, error) {
	userID := helpers.GetUserID(ctx)
	if userID == "" {
		logging.Log.Warn("failed to get userID from token")
		return dto.CommentResponse{}, constants.ErrGetIDFromToken
	}

	comment := model.Comment{
		ID:         uuid.New(),
		ProposalID: req.ProposalID,
		UserID:     uuid.MustParse(userID),
		Content:    req.Content,
	}

	if err := cs.commentRepo.CreateComment(ctx, nil, comment); err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_CREATE_COMMENT)
		return dto.CommentResponse{}, constants.ErrCreateComment
	}

	logging.Log.WithFields(map[string]interface{}{
		"user_id": userID,
		"comment": req.Content,
	}).Info(constants.MESSAGE_SUCCESS_CREATE_COMMENT)

	user, _, err := cs.userRepo.GetUserByID(ctx, nil, userID)
	if err != nil {
		logging.Log.WithError(err).Error("failed to get user info")
	}

	proposal, _, err := cs.proposalRepo.GetProposalByID(ctx, nil, req.ProposalID.String())
	if err != nil {
		logging.Log.WithError(err).Error("failed to get proposal info")
	}

	return dto.CommentResponse{
		ID:      comment.ID,
		Content: comment.Content,
		User: dto.UserCommentResponse{
			ID:   uuid.MustParse(userID),
			Name: user.Name,
		},
		Proposal: dto.ProposalCommentResponse{
			ID:    proposal.ID,
			Title: proposal.Title,
		},
	}, nil
}

func (cs *CommentService) DeleteComment(ctx context.Context, req dto.DeleteCommentRequest) (dto.CommentResponse, error) {
	userID := helpers.GetUserID(ctx)
	if userID == "" {
		logging.Log.Warn("failed to get userID from token")
		return dto.CommentResponse{}, constants.ErrGetIDFromToken
	}

	// Pastikan comment ada dulu (opsional, untuk validasi)
	comment, _, err := cs.commentRepo.GetCommentByID(ctx, nil, req.CommentID)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_GET_COMMENT)
		return dto.CommentResponse{}, constants.ErrGetCommentByID
	}

	// Pastikan yang menghapus adalah pemilik comment (opsional, tapi bagus untuk keamanan)
	if comment.UserID.String() != userID {
		logging.Log.Warnf("unauthorized delete attempt by user %s on comment %s", userID, req.CommentID)
		return dto.CommentResponse{}, constants.ErrDeniedAccess
	}

	// Hapus comment
	if err := cs.commentRepo.DeleteComment(ctx, nil, req.CommentID); err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_DELETE_COMMENT)
		return dto.CommentResponse{}, constants.ErrDeleteComment
	}

	// Logging sukses
	logging.Log.WithFields(map[string]interface{}{
		"user_id":    userID,
		"comment_id": req.CommentID,
	}).Info(constants.MESSAGE_SUCCESS_DELETE_COMMENT)

	return dto.CommentResponse{
		ID: comment.ID,
		User: dto.UserCommentResponse{
			ID:   comment.UserID,
			Name: "", // optional: bisa ambil nama user dari repo kalau mau tampil
		},
		Content: comment.Content,
	}, nil
}

