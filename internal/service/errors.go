package service

import "errors"

var (
	ErrStudentInDb = errors.New("student already in the database")
	ErrStudentNotFound = errors.New("student not in the database")
)