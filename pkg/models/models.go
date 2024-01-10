package models

import (
	"errors"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type Users struct {
	UserID   int
	Gender   string
	Username string
	Password string
	FullName string
}

type Boards struct {
	BoardID    int
	BoardName  string
	CreateDate string
	Task       []Tasks
	UserID     int
}

type Tasks struct {
	TaskID      int
	TaskName    string
	Priority    string
	Description string
}

type TasksBoards struct {
	TaskID  int
	BoardID int
	Status  string
}
