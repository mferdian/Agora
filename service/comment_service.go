package service

import (
	"Agora/model"
	"Agora/repository"
	"Agora/websocket"
	"github.com/google/uuid"
	"time"
)

type ICommentService interface {
	CreateComment(userID string, req model.Comment) (model.Comment, error)
}

type commentService struct {
	repo repository.ICommentRepository
}

func NewCommentService(repo repository.ICommentRepository) ICommentService {
	return &commentService{repo: repo}
}

func (s *commentService) CreateComment(userID string, req model.Comment) (model.Comment, error) {
	req.ID = uuid.New()
	req.UserID, _ = uuid.Parse(userID)
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	comment, err := s.repo.Create(req)
	if err != nil {
		return model.Comment{}, err
	}

	// Kirim notifikasi ke semua client WebSocket
	websocket.WsHub.BroadcastEvent("comment_created", comment)

	return comment, nil
}
