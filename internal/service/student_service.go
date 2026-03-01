package service

import (
	"sort"
	"suai-queue/pkg/student"
	"sync"
)

type StudentService struct {
	students []student.Student
	mutex    sync.Mutex
}

func NewStudentService() *StudentService {
	return &StudentService{
		students: make([]student.Student, 0, 10),
		mutex:    sync.Mutex{},
	}
}

func (serv *StudentService) GetStudents() []student.Student {
	serv.mutex.Lock()
	defer serv.mutex.Unlock()

	DbCopy := make([]student.Student, len(serv.students), cap(serv.students))
	copy(DbCopy, serv.students)
	return DbCopy
}

func (serv *StudentService) Exists(id int64) bool {
	serv.mutex.Lock()
	defer serv.mutex.Unlock()

	return serv.exists(id)
}

func (serv *StudentService) exists(id int64) bool {
	for _, st := range serv.students {
		if st.ID == id {
			return true
		}
	}
	return false
}

func (serv *StudentService) Insert(s *student.Student) error {
	serv.mutex.Lock()
	defer serv.mutex.Unlock()

	if serv.exists(s.ID) {
		return ErrStudentInDb
	}

	i := sort.Search(len(serv.students), func(i int) bool {
		return serv.students[i].ID >= s.ID
	})

	serv.students = append(serv.students, student.Student{})
	copy(serv.students[i+1:], serv.students[i:])
	serv.students[i] = *s

	return nil
}

func (serv *StudentService) Remove(ID int64) error {
	serv.mutex.Lock()
	defer serv.mutex.Unlock()

	i := sort.Search(len(serv.students), func(i int) bool {
		return serv.students[i].ID >= ID
	})

	if i >= len(serv.students) || serv.students[i].ID != ID {
		return ErrStudentNotFound
	}

	serv.students = append(serv.students[:i], serv.students[i+1:]...)
	return nil
}

func (serv *StudentService) GetName(id int64) string {
	serv.mutex.Lock()
	defer serv.mutex.Unlock()

	for _, s := range serv.students {
		if s.ID == id {
			return s.Name
		}
	}

	return ""
}

func (serv *StudentService) UpdateName(id int64, newName string) error {
	serv.mutex.Lock()
	defer serv.mutex.Unlock()

	for i := range serv.students {
		if serv.students[i].ID == id {
			serv.students[i].Name = newName
			return nil
		}
	}

	return ErrStudentNotFound
}
