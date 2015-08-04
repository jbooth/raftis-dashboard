package main

import (
	"fmt"
	//"github.com/coreos/go-etcd/etcd"
	"log"
	"net/http"
)

func main() {
	//client := etcd.NewClient([]string{"http://127.0.0.1:4001"})
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world")
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
