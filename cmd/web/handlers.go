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

	board, err := app.boards.GetCurrentBoard()
	if err != nil {
		app.serverError(w, err)
		return
	}

	fmt.Fprintf(w, "%v\n", board)

}

func (app *application) showBoard(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	s, err := app.boards.GetBoard(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return

	}

	files := []string{
		"./ui/html/show.board.html",
		"./ui/html/base.layout.html",
		"./ui/html/footer.partial.html",
	}

	// Парсинг файлов шаблонов...
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// А затем выполняем их. Обратите внимание на передачу заметки с данными
	// (структура models.Boards) в качестве последнего параметра.
	err = ts.Execute(w, s)
	if err != nil {
		app.serverError(w, err)
	}

	// Отображаем весь вывод на странице.
	fmt.Fprintf(w, "%v", s)
}

func (app *application) createBoard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed) // Используем помощник clientError()
		return
	}
	id, err := app.boards.InsertBoard("test_board", 1, []int{1})
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/kanban?id=%d", id), http.StatusSeeOther)

	w.Write([]byte("Создание новой kanban..."))
}
