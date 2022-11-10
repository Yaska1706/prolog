package main

import (
	"log"

	"github.com/yaska1706/prolog/internal/server"
)

func main() {
	srvr := server.NewHttpServer(":8080")
	log.Fatal(srvr.ListenAndServe())
}
