package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type User struct {
	ID        int
	Login     string
	Gender    string
	BirthYear int
	KanbanID  int
}

type Kanban struct {
	ID           int
	TaskID       int
	Expires      time.Time
	CreationDate time.Time
}

type Task struct {
	ID           int
	UserID       int
	Name         string
	Priotity     string
	CreationDate time.Time
}
