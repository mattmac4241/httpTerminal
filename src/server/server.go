package main

import (
	"commands"
	"net/http"
)

type Message struct {
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		//from := r.FormValue("From")
		body := r.FormValue("Body")
		w.Write([]byte(commands.EvalCommand(body)))
	}
}
