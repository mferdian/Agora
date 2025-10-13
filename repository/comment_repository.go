package repository

import (
	"Agora/model"
	"context"

	"gorm.io/gorm"
)

type ICommentRepository interface {
	CreateComment(ctx context.Context, tx *gorm.DB, comment model.Comment) error
	DeleteComment(ctx context.Context, tx *gorm.DB, id string) error
	GetCommentByID(ctx context.Context, tx *gorm.DB, id string) (model.Comment, bool, error)
}

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) ICommentRepository {
	return &CommentRepository{db: db}
}

func (cr *CommentRepository)CreateComment(ctx context.Context, tx *gorm.DB, comment model.Comment) error {
	if tx == nil {
		tx = cr.db
	}

	return tx.WithContext(ctx).Create(&comment).Error
}

func (cr *CommentRepository) GetCommentByID(ctx context.Context, tx *gorm.DB, id string) (model.Comment, bool, error) {
	if tx == nil {
		tx = cr.db
	}

	var comment model.Comment
	if err := tx.WithContext(ctx).Where("id = ?", id).Take(&comment).Error; err != nil {
		return model.Comment{}, false, err
	}

	return comment, true, nil
}


func (cr *CommentRepository) DeleteComment(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = cr.db
	}

	return tx.WithContext(ctx).Where("id = ?", id).Delete(&model.Comment{}).Error
}
