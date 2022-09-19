package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

// templ represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, nil)
}
func main() {
	// handlefunc maps a path : / to a function: here
	template := &templateHandler{filename: "chat.html"}

	// create a Room
	r := newRoom()

	http.HandleFunc("/", template.ServeHTTP)
	http.HandleFunc("/room", r.ServeHTTP)
	// get the room going : room is going to run in a separate subroutine
	go r.run()

	//starting the web server on the main thread
	if err := http.ListenAndServe(":8080", nil); err != nil {
		// print error
		log.Fatal("ListenAndServe:", err)
	}
}
