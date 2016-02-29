package main

import (
	"commands"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body := r.FormValue("Body")
		w.Write([]byte(commands.EvalCommand(body)))
	}
}
