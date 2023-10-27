package main

import (
	"errors"
	"fmt"
	"html/template"
	"mephi/kanban/pkg/models"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	files := []string{
		"./ui/html/kanban_desk.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) showKanban(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	s, err := app.kanban.GetKanban(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Отображаем весь вывод на странице.
	fmt.Fprintf(w, "%v", s)
}

func (app *application) createKanban(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed) // Используем помощник clientError()
		return
	}
	taskId := 1
	id, err := app.kanban.InsertKanban(taskId)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/kanban?id=%d", id), http.StatusSeeOther)

	w.Write([]byte("Создание новой kanban..."))
}
