package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"mephi/kanban/pkg/models"
	"mephi/kanban/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog     *log.Logger
	infoLog      *log.Logger
	boards       *mysql.BoardsModel
	tasks        *mysql.TasksModel
	tasks_boards *mysql.TasksBoardsModel
	users        *mysql.UsersModel
	cache        map[string]*models.Users
}

func main() {
	addr := flag.String("addr", ":4000", "Сетевой адрес веб-сервера")
	dsn := flag.String("dsn", "cat:1111@/kanban?parseTime=true", "Название MySQL источника данных")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		errorLog:     errorLog,
		infoLog:      infoLog,
		boards:       &mysql.BoardsModel{DB: db},
		tasks:        &mysql.TasksModel{DB: db},
		users:        &mysql.UsersModel{DB: db},
		tasks_boards: &mysql.TasksBoardsModel{DB: db},
		cache:        make(map[string]*models.Users),
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Запуск сервера на %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
