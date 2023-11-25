package mysql

import (
	"database/sql"
	"errors"
	"mephi/kanban/pkg/models"
	"time"
)

type BoardsModel struct {
	DB *sql.DB
}

func (m *BoardsModel) InsertBoard(boardName string, userID int, tasks []int) (int, error) {
	createDate := time.Now()
	stmt := `INSERT INTO boards (boardName, userID, createDate)
    VALUES(?, ?, UTC_TIMESTAMP())`

	result, err := m.DB.Exec(stmt, boardName, userID, createDate)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	for _, taskID := range tasks {
		query := "INSERT INTO tasks_boards (taskID, boardID) VALUES ($1, $2)"
		_, err = m.DB.Exec(query, taskID, id)
		if err != nil {
			return 0, err
		}
	}

	return int(id), nil
}

func (m *BoardsModel) GetBoard(boardID int) (*models.Boards, error) {
	stmt := `SELECT boardID, boardName, userID, createDate FROM boards
    WHERE boardID = ?`

	row := m.DB.QueryRow(stmt, boardID)

	s := &models.Boards{}

	err := row.Scan(&s.BoardID, &s.BoardName, &s.UserID, &s.CreateDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *BoardsModel) GetUsersBoards(userID int) ([]*models.Boards, error) {
	stmt := `SELECT boards.boardID, boards.boardName, boards.createDate FROM boards
    WHERE boards.userID = ?`

	rows, err := m.DB.Query(stmt, userID)
	defer rows.Close()

	//s := &[]models.Boards{}
	var boards []*models.Boards

	for rows.Next() {
		board := &models.Boards{}
		err := rows.Scan(&board.BoardID, &board.BoardName, &board.CreateDate)
		if err != nil {
			return nil, err
		}
		boards = append(boards, board)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return boards, nil
}

func (m *BoardsModel) GetCurrentBoard() (*models.Boards, error) {
	stmt := `SELECT boards.boardID, boards.boardName, boards.userID, boards.createDate FROM boards
    WHERE DATE_ADD(board.createDate, INTERVAL 2 WEEK) > (CURDATE() + board.createDate)`

	row := m.DB.QueryRow(stmt)

	s := &models.Boards{}

	err := row.Scan(&s.BoardID, &s.BoardName, &s.UserID, &s.CreateDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}
