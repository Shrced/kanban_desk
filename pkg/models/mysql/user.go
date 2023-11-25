package mysql

import (
	"database/sql"
	"errors"
	"mephi/kanban/pkg/models"
)

type UsersModel struct {
	DB *sql.DB
}

func (m *UsersModel) InsertUser(gender string, username string, password string, fullName string) (int, error) {
	stmt := `INSERT INTO users (gender, username, password, fullName)
    VALUES(?, ?, ?, ?)`

	result, err := m.DB.Exec(stmt, gender, username, password, fullName)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *UsersModel) GetUser(UserID int) (*models.Users, error) {
	stmt := `SELECT gender, username, password, fullName FROM users
    WHERE userID = ?`

	row := m.DB.QueryRow(stmt, UserID)

	s := &models.Users{}

	err := row.Scan(&s.Gender, &s.Username, &s.Password, &s.FullName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}