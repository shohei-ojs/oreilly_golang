package main

import (
	"log"
	"net/http"
	"text/template"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ =
				template.Must(template.ParseFiles(filepath.Join("templates",
						t.filename)))
	})
	t.templ.Execute(w, nil)
}

func main() {
	// root
	http.Handle("/", &templateHandler{filename: "chat.html"})
	// Webサーバーを開始します
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}