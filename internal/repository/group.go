package repository

import (
	"errors"
	"suai-queue/internal/domain"

	"gorm.io/gorm"
)

type GroupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) Create(groupNumber string) error {
	newGroup := domain.Group{Number: groupNumber}

	result := r.db.Create(&newGroup)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrGroupAlreadyExists
		}
		return result.Error
	}
	return nil
}

func (r *GroupRepository) GetAll() ([]domain.Group, error) {
	var groups []domain.Group
	err := r.db.Find(&groups).Error
	return groups, err
}

func (r *GroupRepository) Exists(groupNumber string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Group{}).Where("number = ?", groupNumber).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
