package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/go-zoo/bone"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 2222, "port number")

	mux := bone.New()
	mux.GetFunc("/authorize", Auth)
	mux.GetFunc("/auth_callback", Callback)

	http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

func Auth(w http.ResponseWriter, r *http.Request) {

}

func Callback(w http.ResponseWriter, r *http.Request) {

}
