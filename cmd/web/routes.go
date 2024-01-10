package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.authorized(app.home))
	mux.HandleFunc("/kanban", app.authorized(app.showBoard))
	mux.HandleFunc("/task/create", app.authorized(app.createTask))
	mux.HandleFunc("/task", app.authorized(app.showTask))
	mux.HandleFunc("/task/update", app.authorized(app.updateTask))
	mux.HandleFunc("/kanban/kanban_list", app.authorized(app.showBoards))
	mux.HandleFunc("/login", app.Login)
	mux.HandleFunc("/logout", app.authorized(app.Logout))
	mux.HandleFunc("/signup", app.SignUp)
	mux.HandleFunc("/user/profile", app.authorized(app.showProfile))
	mux.HandleFunc("/user/current_desk", app.authorized(app.home))

	fileServer := http.FileServer(http.Dir("/home/kottik/code/kanban/"))
	mux.Handle("/ui/static/", http.StripPrefix("/ui/static", fileServer))

	return mux
}
