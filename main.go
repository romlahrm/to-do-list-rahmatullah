package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	username = "root"
	password = "password"
	portnum = 3306
	dbname = "todolist"
)

func main() {
	fmt.Println("Hello world")

	r := mux.NewRouter()

	type Todo struct {
		Id    int
		Title string
		Done  bool
	}

	type TodoPageData struct {
		PageTitle string
		Todos     []Todo
	}

	conf := username + ":" + password + "@(localhost:" + strconv.Itoa(portnum) + ")/" + dbname + "?parseTime=true"

	fmt.Println(conf)

	db, err := sql.Open("mysql", conf)

	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("no error")
	}

	homeTempl := template.Must(template.ParseFiles("./templates/index.html"))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		todos, err := db.Query("select * from todos")

		defer todos.Close()

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("no error selecting todos")
		}

		var todo []Todo

		for todos.Next() {
			var t Todo
			err := todos.Scan(&t.Id, &t.Title, &t.Done)
			// fmt.Println(err)
			if err == nil {
				fmt.Println(t)
				todo = append(todo, t)
				// fmt.Println(todo)
			} else {
				fmt.Println(err)
				return
			}
		}

		// if r.Method == http.MethodPost{
		// 	title := r.FormValue("todotitle")
		// 	if title != ""{
		// 		_, err := db.Exec(`insert into todos (title, done) values(?, ?)`, title, 0)
		// 		if err != nil {
		// 			fmt.Println(err)
		// 		}else{
		// 			fmt.Println("Todo inserted successfully")
		// 		}
		// 	}

		// }

		fmt.Println(todo)

		data := TodoPageData{
			PageTitle: "To do list",
			Todos:     todo,
		}

		// data = TodoPageData{
		// 	PageTitle: "To do custom",
		// 	Todos: []Todo{
		// 		{1, "test", true},
		// 		{2, "test", false},
		// 	},
		// }

		// for _, t := range todo {
		// 	fmt.Println(t.title)
		// }

		fmt.Println(data)
		er := homeTempl.Execute(w, data)

		if er != nil {
			fmt.Println(er)
		}
	})

	r.HandleFunc("/remove-todo/{id}", func(w http.ResponseWriter, r *http.Request){
		id := mux.Vars(r)["id"]
		_, err := db.Exec(`Delete from todos where id = ?`, id)

		if err != nil{
			fmt.Println(err)
		}else{
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

	})

	r.HandleFunc("/add-todo", func(w http.ResponseWriter, r *http.Request){
		if r.Method == http.MethodPost {
			title := r.FormValue("todotitle")
			_, err = db.Exec("insert into todos(title) values(?)", title)
			if err != nil{
				fmt.Println(err)
			}else {
				fmt.Println("no errors")
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		}
	})

	// r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
	// 	vars := mux.Vars(r)
	// 	title := vars["title"]
	// 	pages := vars["page"]

	// 	fmt.Fprintf(w, "You have requested the book %s on page %s\n", title, pages)
	// })

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
	// 	fmt.Fprintf(w, "Welcome to my website!!")
	// })

	// fs := http.FileServer(http.Dir("static/"))

	// http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	r.
		PathPrefix("/assets/").
		Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("."+"/assets/"))))

	http.ListenAndServe(":8090", r)
}
