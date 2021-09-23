package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

// templ means a template file
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP serves HTTP requests
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "address of the application")
	flag.Parse()
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	// rooting
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	// start chat room
	go r.run()
	// start WEB server
	log.Println("Starting web server. port:", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
