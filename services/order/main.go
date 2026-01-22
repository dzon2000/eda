package main

import (
	"log"
	"net/http"

	"github.com/dzon2000/eda/order/internal/web"
)

func main() {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: web.NewHandler().Router(),
	}

	log.Fatal(srv.ListenAndServe())
}
