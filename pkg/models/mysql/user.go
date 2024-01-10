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
	stmt := `SELECT gender, username, password, fullName, userID FROM users
    WHERE userID = ?`

	row := m.DB.QueryRow(stmt, UserID)

	s := &models.Users{}

	err := row.Scan(&s.Gender, &s.Username, &s.Password, &s.FullName, &s.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *UsersModel) Login(Username, Password string) (*models.Users, error) {
	stmt := `SELECT gender, username, fullName, userID FROM users
    WHERE password = ? AND username = ?`

	row := m.DB.QueryRow(stmt, Password, Username)

	user := &models.Users{}

	err := row.Scan(&user.Gender, &user.Username, &user.FullName, &user.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return user, nil
}
