package service

import (
	"sort"
	"suai-queue/pkg/student"
	"sync"
)

type StudentService struct {
	students []student.Student
	mutex sync.Mutex
}

func (serv *StudentService) Exists(s student.Student) bool {
    serv.mutex.Lock()
    defer serv.mutex.Unlock()

    for _, st := range serv.students {
        if st.TelegramLogin == s.TelegramLogin {
            return true
        }
    }
    return false
}

func (serv *StudentService) Insert(s student.Student) error {
    serv.mutex.Lock()
    defer serv.mutex.Unlock()

    if serv.Exists(s) {
        return ErrStudentInDb
    }

    i := sort.Search(len(serv.students), func(i int) bool {
        return serv.students[i].TelegramLogin >= s.TelegramLogin
    })

    serv.students = append(serv.students, student.Student{})
    copy(serv.students[i+1:], serv.students[i:])
    serv.students[i] = s

    return nil
}

func (serv *StudentService) Remove(telegramLogin string) error {
    serv.mutex.Lock()
    defer serv.mutex.Unlock()

    i := sort.Search(len(serv.students), func(i int) bool {
        return serv.students[i].TelegramLogin >= telegramLogin
    })

    if i >= len(serv.students) || serv.students[i].TelegramLogin != telegramLogin {
        return ErrStudentNotFound
    }

    serv.students = append(serv.students[:i], serv.students[i+1:]...)
    return nil
}