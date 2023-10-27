package mysql

import (
	"database/sql"
	"errors"
	"mephi/kanban/pkg/models"
	"time"
)

// Kanban - Определяем тип который обертывает пул подключения sql.DB
type KanbanModel struct {
	DB *sql.DB
}

// InsertUser - Метод для создания новой user в базе дынных.
func (m *KanbanModel) InsertUser(BirthYear, KanbanID int, Login, Gender string) (int, error) {
	stmt := `INSERT INTO users (birth_year, kanban_id, login, gender)
    VALUES(?, ?, ?, ?)`

	result, err := m.DB.Exec(stmt, BirthYear, KanbanID, Login, Gender)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// InsertTask - Метод для создания новой task в базе дынных.
func (m *KanbanModel) InsertTask(UserID int, Name string) (int, error) {
	curTime := time.Now()
	stmt := `INSERT INTO tasks (user_id, name, creation_date)
    VALUES(?, ?, UTC_TIMESTAMP())`

	result, err := m.DB.Exec(stmt, UserID, Name, curTime)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// InsertKanban - Метод для создания новой kanban в базе дынных.
func (m *KanbanModel) InsertKanban(TaskID int) (int, error) {
	curTime := time.Now()
	stmt := `INSERT INTO kanban (task_id, creation_date)
    VALUES(?, UTC_TIMESTAMP())`

	result, err := m.DB.Exec(stmt, TaskID, curTime)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// GetTask - Метод для возвращения данных task по идентификатору ID.
func (m *KanbanModel) GetTask(id int) (*models.Task, error) {
	stmt := `SELECT id, user_id, name, creation_date FROM tasks
    WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &models.Task{}

	err := row.Scan(&s.ID, &s.UserID, &s.Name, &s.CreationDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

// GetUser - Метод для возвращения данных пользователя по его идентификатору ID.
func (m *KanbanModel) GetUser(id int) (*models.User, error) {
	stmt := `SELECT id, login, gender, birth_year, kanban_id FROM users
    WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &models.User{}

	err := row.Scan(&s.ID, &s.Login, &s.Gender, &s.BirthYear, &s.KanbanID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

// GetKanban - Метод для возвращения данных kanban desk по идентификатору ID.
func (m *KanbanModel) GetKanban(id int) (*models.Kanban, error) {
	stmt := `SELECT id, task_id, creation_date FROM kanban
    WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &models.Kanban{}

	err := row.Scan(&s.ID, &s.TaskID, &s.CreationDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

// Latest - Метод возвращает 10 наиболее часто используемые заметки.
func (m *KanbanModel) Latest() ([]*models.Task, error) {
	return nil, nil
}
