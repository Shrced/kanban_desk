package mysql

import (
	"database/sql"
	"errors"
	"mephi/kanban/pkg/models"
)

type TasksModel struct {
	DB *sql.DB
}

type TasksBoardsModel struct {
	DB *sql.DB
}

func (m *TasksBoardsModel) InsertTaskBoards(TaskId, BoardId int, status string) error {
	query := "INSERT INTO tasks_boards (taskID, boardID, status) VALUES (?, ?, ?)"
	_, err := m.DB.Exec(query, TaskId, BoardId, status)
	if err != nil {
		return err
	}

	return nil
}

func (m *TasksModel) InsertTask(taskName string, BoardID int, status string, description string, priority string) (int, error) {
	stmt := `INSERT INTO tasks (taskName, description, priority)
    VALUES(?, ?, ?)`

	result, err := m.DB.Exec(stmt, taskName, description, priority)
	if err != nil {
		return 0, err
	}

	TaskID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	query := "INSERT INTO tasks_boards (taskID, boardID, status) VALUES (?, ?, ?)"
	_, err2 := m.DB.Exec(query, TaskID, BoardID, status)
	if err2 != nil {
		return 0, err2
	}

	return int(TaskID), nil
}

func (m *TasksModel) UpdateTask(userID, taskID int, taskName string, status string, description string, priority string) (int, error) {
	stmt := `UPDATE tasks SET taskName = ?, description = ?, priority = ? WHERE taskID = ?`

	_, err := m.DB.Exec(stmt, taskName, description, priority, taskID)
	if err != nil {
		return 0, err
	}

	stmt2 := `SELECT boards.boardID FROM boards
    WHERE DATE_ADD(boards.createDate, INTERVAL 2 WEEK) > CURDATE() AND boards.userID = ?`
	s := &models.Boards{}
	err3 := m.DB.QueryRow(stmt2, userID).Scan(&s.BoardID)
	if err3 != nil {
		return 0, err3
	}

	query := "UPDATE tasks_boards SET status = ? WHERE boardID = ? and taskID = ?"
	_, err2 := m.DB.Exec(query, status, s.BoardID, taskID)
	if err2 != nil {
		return 0, err2
	}

	return taskID, nil
}

func (m *TasksModel) GetTask(taskID int) (*models.Tasks, error) {
	stmt := `SELECT taskID, taskName, description, priority FROM tasks
    WHERE taskID = ?`

	row := m.DB.QueryRow(stmt, taskID)

	s := &models.Tasks{}

	err := row.Scan(&s.TaskID, &s.TaskName, &s.Description, &s.Priority)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *TasksModel) GetBoardsTasks(boardID int) ([]*models.Tasks, []*models.Tasks, []*models.Tasks, error) {
	stmt1 := `SELECT tasks.taskID, tasks.taskName, tasks.description, tasks.priority FROM tasks
    JOIN tasks_boards ON tasks.taskID  = tasks_boards.taskID
    WHERE tasks_boards.boardID = ? and tasks_boards.status = 'in progress'`

	stmt2 := `SELECT tasks.taskID, tasks.taskName, tasks.description, tasks.priority FROM tasks 
    JOIN tasks_boards ON tasks.taskID  = tasks_boards.taskID
    WHERE tasks_boards.boardID = ? and tasks_boards.status = 'done'`

	stmt3 := `SELECT tasks.taskID, tasks.taskName, tasks.description, tasks.priority FROM tasks 
    JOIN tasks_boards ON tasks.taskID  = tasks_boards.taskID
    WHERE tasks_boards.boardID = ? and tasks_boards.status = 'to do'`

	rows1, _ := m.DB.Query(stmt1, boardID)
	defer rows1.Close()

	rows2, _ := m.DB.Query(stmt2, boardID)
	defer rows2.Close()

	rows3, _ := m.DB.Query(stmt3, boardID)
	defer rows3.Close()

	// in progress
	var tasks_Progress []*models.Tasks

	for rows1.Next() {
		task := &models.Tasks{}
		err := rows1.Scan(&task.TaskID, &task.TaskName, &task.Description, &task.Priority)
		if err != nil {
			return nil, nil, nil, err
		}

		tasks_Progress = append(tasks_Progress, task)
	}

	// done
	var tasksDone []*models.Tasks

	for rows2.Next() {
		task := &models.Tasks{}
		err := rows2.Scan(&task.TaskID, &task.TaskName, &task.Description, &task.Priority)
		if err != nil {
			return nil, nil, nil, err
		}

		tasksDone = append(tasksDone, task)
	}
	//to do
	var tasksToDo []*models.Tasks

	for rows3.Next() {
		task := &models.Tasks{}
		err := rows3.Scan(&task.TaskID, &task.TaskName, &task.Description, &task.Priority)
		if err != nil {
			return nil, nil, nil, err
		}

		tasksToDo = append(tasksToDo, task)
	}

	return tasksToDo, tasks_Progress, tasksDone, nil
}
