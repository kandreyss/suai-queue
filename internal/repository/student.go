package repository

import (
	"suai-queue/internal/domain"

	"gorm.io/gorm"
)

type StudentRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (s *StudentRepository) Create(st *domain.Student) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var g domain.Group
		if err := tx.Where("number = ?", st.Group).First(&g).Error; err != nil {
			return ErrGroupNotFound
		}

		if err := tx.Create(st).Error; err != nil {
			return err
		}

		return tx.Model(&domain.Group{}).Where("number = ?", st.Group).
			Update("counter", gorm.Expr("counter + 1")).Error
	})
}

func (s *StudentRepository) GetByTGID(tgID int64) (*domain.Student, error) {
	var st domain.Student

	err := s.db.Where("tg_id = ?", tgID).First(&st).Error
	if err != nil {
		return nil, err
	}
	return &st, nil
}

func (s *StudentRepository) Exists(tgID int64) bool {
	var count int64
	s.db.Model(&domain.Student{}).Where("tg_id = ?", tgID).Count(&count)
	return count > 0
}

func (s *StudentRepository) Update(st *domain.Student) error {
	return s.db.Save(st).Error
}

func (s *StudentRepository) Delete(st *domain.Student) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		groupNumber := st.Group

		if err := tx.Delete(st).Error; err != nil {
			return err
		}

		if groupNumber != "" {
			return tx.Model(&domain.Group{}).Where("number = ?", groupNumber).
				Update("counter", gorm.Expr("counter - 1")).Error
		}
		return nil
	})
}

func (s *StudentRepository) GetName(tgID int64) string {
	var st domain.Student

	if err := s.db.Select("name").Where("tg_id = ?", tgID).First(&st).Error; err != nil {
		return ""
	}
	return st.Name
}

func (s *StudentRepository) UpdateName(tgID int64, name string) error {
	return s.db.Model(&domain.Student{}).Where("tg_id = ?", tgID).Update("name", name).Error
}

func (s *StudentRepository) UpdateGroup(tgID int64, groupNumber string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var g domain.Group
		if err := tx.Where("number = ?", groupNumber).First(&g).Error; err != nil {
			return ErrGroupNotFound
		}

		var st domain.Student
		if err := tx.Where("tg_id = ?", tgID).First(&st).Error; err != nil {
			return err
		}
		oldGroup := st.Group

		if err := tx.Model(&domain.Student{}).Where("tg_id = ?", tgID).
			Update("group_number", groupNumber).Error; err != nil {
			return err
		}

		if oldGroup != "" {
			if err := tx.Model(&domain.Group{}).Where("number = ?", oldGroup).
				Update("counter", gorm.Expr("counter - 1")).Error; err != nil {
				return err
			}
		}

		return tx.Model(&domain.Group{}).Where("number = ?", groupNumber).
			Update("counter", gorm.Expr("counter + 1")).Error
	})
}
