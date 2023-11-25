package models

import (
	"errors"
	"time"
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
	CreateDate time.Time
	Task       []Tasks
	UserID     int
}

type Tasks struct {
	TaskID   int
	TaskName string
	Status   string
	Board    []Boards
}

type TasksBoards struct {
	TaskID  int
	BoardID int
}
