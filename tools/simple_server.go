package main

import "net/http"

func main() {
	println("Serving on http://localhost:8080/index.html")
	panic(http.ListenAndServe(":8080", http.FileServer(http.Dir("http"))))
}
