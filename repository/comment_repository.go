package repository

import (
	"Agora/model"
	"gorm.io/gorm"
)

type ICommentRepository interface {
	Create(comment model.Comment) (model.Comment, error)
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) ICommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(comment model.Comment) (model.Comment, error) {
	if err := r.db.Create(&comment).Error; err != nil {
		return model.Comment{}, err
	}
	return comment, nil
}
