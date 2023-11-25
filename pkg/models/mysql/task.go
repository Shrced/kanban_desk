package mysql

import (
	"database/sql"
	"errors"
	"mephi/kanban/pkg/models"
)

type TasksModel struct {
	DB *sql.DB
}

func (m *TasksModel) InsertTask(taskName string, status string, boards []int) (int, error) {
	stmt := `INSERT INTO tasks (taskName, status)
    VALUES(?, ?)`

	result, err := m.DB.Exec(stmt, taskName, status)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	for _, boardID := range boards {
		query := "INSERT INTO tasks_boards (taskID, boardID) VALUES ($1, $2)"
		_, err = m.DB.Exec(query, id, boardID)
		if err != nil {
			return 0, err
		}
	}

	return int(id), nil
}

func (m *TasksModel) GetTask(taskID int) (*models.Tasks, error) {
	stmt := `SELECT taskID, taskName, status FROM tasks
    WHERE taskID = ?`

	row := m.DB.QueryRow(stmt, taskID)

	s := &models.Tasks{}

	err := row.Scan(&s.TaskID, &s.TaskName, &s.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *BoardsModel) GetBoardsTasks(boardID int) ([]*models.Tasks, error) {
	stmt := `SELECT tasks.taskID, tasks.taskName, tasks.status FROM tasks 
    JOIN tasks_boards ON tasks.taskID  = tasks_boards.taskID
    WHERE tasks_boards.userID = ?`

	rows, err := m.DB.Query(stmt, boardID)
	defer rows.Close()

	//s := &[]models.Boards{}
	var tasks []*models.Tasks

	for rows.Next() {
		task := &models.Tasks{}
		err := rows.Scan(&task.TaskID, &task.TaskName, &task.Status)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return tasks, nil
}
