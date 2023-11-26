package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/kanban", app.showBoard)
	mux.HandleFunc("/task/create", app.createTask)
	mux.HandleFunc("/task", app.showTask)
	mux.HandleFunc("/kanban_list", app.showBoards)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
