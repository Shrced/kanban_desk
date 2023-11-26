package main

import "mephi/kanban/pkg/models"

type templateData struct {
	Boards     *models.Boards
	BoardsList []*models.Boards
	Tasks      *models.Tasks
	TasksList  []*models.Tasks
	Users      *models.Users
}
