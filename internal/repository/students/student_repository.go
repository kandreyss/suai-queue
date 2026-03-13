package students

import (
	"suai-queue/pkg/student"

	"gorm.io/gorm"
)

type StudentRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *StudentRepository {
	return &StudentRepository{db: db}
}
func (s *StudentRepository) Create(student *student.Student) error {
	return s.db.Create(student).Error
}

func (s *StudentRepository) GetByTGID(tgID int64) (*student.Student, error) {
	var st student.Student

	err := s.db.Where("tg_id = ?", tgID).First(&st).Error
	if err != nil {
		return nil, err
	}
	return &st, nil
}

func (s *StudentRepository) Exists(tgID int64) bool {
	var count int64
	s.db.Model(&student.Student{}).Where("tg_id = ?", tgID).Count(&count)
	return count > 0
}

func (s *StudentRepository) Update(student *student.Student) error {
	return s.db.Save(student).Error
}

func (s *StudentRepository) Delete(student *student.Student) error {
	return s.db.Delete(student).Error
}
