package mysql

import (
	"database/sql"
	"errors"
	"mephi/kanban/pkg/models"
)

type BoardsModel struct {
	DB *sql.DB
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
	stmt := `SELECT boards.boardID, boards.boardName, boards.createDate, boards.userID FROM boards
    WHERE boards.userID = ?`

	rows, err := m.DB.Query(stmt, userID)

	defer rows.Close()

	var boards []*models.Boards

	for rows.Next() {
		board := &models.Boards{}

		err := rows.Scan(&board.BoardID, &board.BoardName, &board.CreateDate, &board.UserID)
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

func (m *BoardsModel) GetCurrentBoard(userID int, boardName string) (*models.Boards, error) {
	stmt := `SELECT boards.boardID, boards.boardName, boards.userID, boards.createDate FROM boards
    WHERE DATE_ADD(boards.createDate, INTERVAL 2 WEEK) > CURDATE() AND boards.userID = ?`
	s := &models.Boards{}
	err := m.DB.QueryRow(stmt, userID).Scan(&s.BoardID, &s.BoardName, &s.UserID, &s.CreateDate)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {

			// Вставка новой канбан
			stmt := `INSERT INTO boards (boardName, userID, createDate)
			VALUES(?, ?, DATE(UTC_TIMESTAMP()))`

			result, err := m.DB.Exec(stmt, boardName, userID)

			if err != nil {
				return nil, err
			}

			id, err := result.LastInsertId()

			stmt2 := `SELECT t.taskID, t.status from tasks_boards as t inner join boards as b on b.boardID = t.boardID WHERE t.status != "done" AND b.userID = ?`

			rows, err := m.DB.Query(stmt2, userID)

			defer rows.Close()

			var tasksBoards []*models.TasksBoards

			for rows.Next() {
				tasksBoard := &models.TasksBoards{}

				err := rows.Scan(&tasksBoard.TaskID, &tasksBoard.Status)
				if err != nil {
					return nil, err
				}
				tasksBoards = append(tasksBoards, tasksBoard)
			}

			for _, taskBoard := range tasksBoards {
				query := "INSERT INTO tasks_boards (taskID, boardID, status) VALUES (?, ?, ?)"
				_, err = m.DB.Exec(query, taskBoard.TaskID, id, taskBoard.Status)
				if err != nil {
					return s, err
				}
			}

			stmt5 := `SELECT boards.boardID, boards.boardName, boards.userID, boards.createDate FROM boards
			WHERE DATE_ADD(boards.createDate, INTERVAL 2 WEEK) > CURDATE() AND boards.userID = ?`
			s := &models.Boards{}
			err5 := m.DB.QueryRow(stmt5, userID).Scan(&s.BoardID, &s.BoardName, &s.UserID, &s.CreateDate)
			if err5 != nil {
				return s, nil
			}

		} else {
			return nil, err
		}
	}

	return s, nil

}

func (m *BoardsModel) GetCurrentBoardID(userID int) (int, error) {
	stmt := `SELECT boards.boardID FROM boards
    WHERE DATE_ADD(boards.createDate, INTERVAL 2 WEEK) > CURDATE() AND boards.userID = ?`
	s := &models.Boards{}
	err := m.DB.QueryRow(stmt, userID).Scan(&s.BoardID)
	if err == nil {
		return 0, nil

	}
	return s.BoardID, nil
}
