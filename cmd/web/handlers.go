package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"log"
	"mephi/kanban/pkg/models"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var user_id_global int

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

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}

		Username := r.FormValue("username")
		Password := r.FormValue("password")

		if Username == "" || Password == "" {
			app.Login(w, r)
			return
		}

		// hash := md5.Sum([]byte(Password))
		// hashedPass := hex.EncodeToString(hash[:])

		user, err := app.users.Login(Username, Password)
		if err != nil {
			app.Login(w, r)
			return
		}
		time64 := time.Now().Unix()
		timeInt := string(time64)
		token := Username + Password + timeInt
		hashToken := md5.Sum([]byte(token))
		hashedToken := hex.EncodeToString(hashToken[:])
		app.cache[hashedToken] = user
		livingTime := 60 * time.Minute
		expiration := time.Now().Add(livingTime)
		//кука будет жить 1 час
		cookie := http.Cookie{Name: "token", Value: url.QueryEscape(hashedToken), Expires: expiration}
		http.SetCookie(w, &cookie)

		user_id_global = user.UserID

		http.Redirect(w, r, fmt.Sprintf("/"), http.StatusSeeOther)

	} else {
		type answer struct {
			Message string
		}
		message := "Opps"
		data := answer{message}
		// data := &templateData{Boards: board}
		files := []string{
			"/home/kottik/code/kanban/ui/html/login.html",
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

func (a *application) authorized(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		token, err := readCookie("token", r)
		if err != nil {
			http.Redirect(rw, r, "/login", http.StatusSeeOther)
			return
		}
		if _, ok := a.cache[token]; !ok {
			http.Redirect(rw, r, "/login", http.StatusSeeOther)
			return
		}
		next(rw, r)
	}
}

func readCookie(name string, r *http.Request) (value string, err error) {
	if name == "" {
		return value, errors.New("you are trying to read empty cookie")
	}
	cookie, err := r.Cookie(name)
	if err != nil {
		return value, err
	}
	str := cookie.Value
	value, _ = url.QueryUnescape(str)
	return value, err
}

func (app *application) Logout(rw http.ResponseWriter, r *http.Request) {
	for _, v := range r.Cookies() {
		c := http.Cookie{
			Name:   v.Name,
			MaxAge: -1}
		http.SetCookie(rw, &c)
	}
	http.Redirect(rw, r, "/login", http.StatusSeeOther)
}

func (app *application) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}

		Username := r.FormValue("username")
		Password := r.FormValue("password")
		Password2 := r.FormValue("password2")
		FullName := r.FormValue("FullName")
		Gender := r.FormValue("gender")

		if Username == "" || Password == "" || FullName == "" || Gender == "" {
			app.SignUp(w, r)
			return
		}

		if Password != Password2 {
			app.SignUp(w, r)
			return
		}

		// hash := md5.Sum([]byte(Password))
		// hashedPass := hex.EncodeToString(hash[:])

		user_id, err := app.users.InsertUser(Gender, Username, Password, FullName)
		if err != nil {
			app.serverError(w, err)
			return
		}

		user_id_global = user_id

		http.Redirect(w, r, fmt.Sprintf("/login"), http.StatusSeeOther)

	} else {
		type answer struct {
			Message string
		}
		message := "Opps"
		data := answer{message}
		// data := &templateData{Boards: board}
		files := []string{
			"/home/kottik/code/kanban/ui/html/signup.html",
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
