package main

import "mephi/kanban/pkg/models"

type templateData struct {
	Boards     *models.Boards
	BoardsList []*models.Boards
	Tasks      *models.Tasks
	TasksList1 []*models.Tasks
	TasksList2 []*models.Tasks
	TasksList3 []*models.Tasks
	Users      *models.Users
}
