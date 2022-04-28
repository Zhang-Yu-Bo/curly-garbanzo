package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Zhang-Yu-Bo/curly-garbanzo/router"
)

func main() {
	mPort := os.Getenv("PORT")
	if mPort == "" {
		mPort = "80"
	}

	mRouter := router.NewRouter()
	mServer := &http.Server{
		Handler:      mRouter,
		Addr:         ":" + mPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("server run on port:", mPort)
	log.Fatal(mServer.ListenAndServe())
}
