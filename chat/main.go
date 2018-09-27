package main

import (
	"log"
	"net/http"
	"text/template"
	"path/filepath"
	"sync"
	"flag"
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
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", "80", "アプリケーションのアドレス")
	flag.Parse() // フラグを解釈
	r := newRoom()
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	// チャットを開始
	go r.run()
	// Webサーバーを開始
	log.Println("Webサーバーを開始します。 ポート : ", *addr)
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}