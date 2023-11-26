package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
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
	tasks_list, err := app.tasks.GetBoardsTasks(board.BoardID)

	data := &templateData{Boards: board, TasksList: tasks_list}

	files := []string{
		"/home/kottik/code/kanban/ui/html/kanban_desk.html",
		"/home/kottik/code/kanban/ui/html/base.layout.html",
		"/home/kottik/code/kanban/ui/html/footer.partial.html",
		"/home/kottik/code/kanban/ui/html/header.partial.html",
	}

	// Парсинг файлов шаблонов...
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// А затем выполняем их. Обратите внимание на передачу заметки с данными
	// (структура models.Boards) в качестве последнего параметра.
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
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
	tasks_list, err := app.tasks.GetBoardsTasks(id)

	data := &templateData{Boards: s, TasksList: tasks_list}

	files := []string{
		"/home/kottik/code/kanban/ui/html/kanban_desk.html",
		"/home/kottik/code/kanban/ui/html/base.layout.html",
		"/home/kottik/code/kanban/ui/html/footer.partial.html",
		"/home/kottik/code/kanban/ui/html/header.partial.html",
	}

	// Парсинг файлов шаблонов...
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// А затем выполняем их. Обратите внимание на передачу заметки с данными
	// (структура models.Boards) в качестве последнего параметра.
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}

}

func (app *application) showBoards(w http.ResponseWriter, r *http.Request) {
	s, err := app.boards.GetUsersBoards(1)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return

	}

	data := &templateData{BoardsList: s}

	files := []string{
		"/home/kottik/code/kanban/ui/html/boards_list.html",
		"/home/kottik/code/kanban/ui/html/base.layout.html",
		"/home/kottik/code/kanban/ui/html/footer.partial.html",
		"/home/kottik/code/kanban/ui/html/header.partial.html",
	}

	// Парсинг файлов шаблонов...
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// А затем выполняем их. Обратите внимание на передачу заметки с данными
	// (структура models.Boards) в качестве последнего параметра.
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

// func (app *application) createBoard(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		w.Header().Set("Allow", http.MethodPost)
// 		app.clientError(w, http.StatusMethodNotAllowed) // Используем помощник clientError()
// 		return
// 	}
// 	id, err := app.boards.InsertBoard("test_board", 1, []int{1})
// 	if err != nil {
// 		app.serverError(w, err)
// 		return
// 	}
// 	http.Redirect(w, r, fmt.Sprintf("/kanban?id=%d", id), http.StatusSeeOther)

// 	w.Write([]byte("Создание новой kanban..."))
// }

func (app *application) createTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}

		TaskName := r.FormValue("task_name")
		Status := r.FormValue("status")
		BoardID, err := strconv.Atoi(r.FormValue("board_id"))
		if err != nil {
			panic(err)
		}

		id, err := app.tasks.InsertTask(TaskName, Status, BoardID)
		// app.tasks_boards.InsertTaskBoards(id, BoardID)

		if err != nil {
			app.serverError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/task?id=%d", id), http.StatusSeeOther)

		w.Write([]byte("Создание новой задачи..."))

	} else {
		board, err := app.boards.GetCurrentBoard()
		if err != nil {
			app.serverError(w, err)
			return
		}

		data := &templateData{Boards: board}
		files := []string{
			"/home/kottik/code/kanban/ui/html/new_task.form.html",
			"/home/kottik/code/kanban/ui/html/base.layout.html",
			"/home/kottik/code/kanban/ui/html/footer.partial.html",
			"/home/kottik/code/kanban/ui/html/header.partial.html",
		}

		// Парсинг файлов шаблонов...
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.serverError(w, err)
			return
		}
		err = ts.Execute(w, data)
		if err != nil {
			app.serverError(w, err)
		}

	}

}

func (app *application) showTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	task, err := app.tasks.GetTask(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return

	}

	data := &templateData{Tasks: task}

	files := []string{
		"/home/kottik/code/kanban/ui/html/task.html",
		"/home/kottik/code/kanban/ui/html/base.layout.html",
		"/home/kottik/code/kanban/ui/html/footer.partial.html",
		"/home/kottik/code/kanban/ui/html/header.partial.html",
	}

	// Парсинг файлов шаблонов...
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// А затем выполняем их. Обратите внимание на передачу заметки с данными
	// (структура models.Boards) в качестве последнего параметра.
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}

}
