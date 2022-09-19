package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	trace "github.com/ash-kamrip/tracer"
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
	// nil here means we aren't passing the data from the request
	// r specified here let's us extract data from http.Request
	t.templ.Execute(w, r)
}
func main() {

	var addr = flag.String("addr", ":8080", "The addr of the application")

	flag.Parse()

	// handlefunc maps a path : / to a function: here
	template := &templateHandler{filename: "chat.html"}

	// create a Room
	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	http.HandleFunc("/", template.ServeHTTP)
	http.HandleFunc("/room", r.ServeHTTP)
	// get the room going : room is going to run in a separate subroutine
	go r.run()

	//starting the web server on the main thread
	log.Println("Starting the web server on ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		// print error
		log.Fatal("ListenAndServe:", err)
	}
}
