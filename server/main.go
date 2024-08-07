package main

import "net/http"

func main() {

	http.HandleFunc("/", HandleHTTPReq)
	http.ListenAndServe(":8000", nil)

}
