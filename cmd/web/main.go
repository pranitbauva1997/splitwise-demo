package main

import (
	"log"

	"github.com/pranitbauva1997/splitwise-demo/pkg"
)

func main() {
	app, err := pkg.Init()
	if err != nil {
		log.Fatalln("couldn't initialize the application:", err)
	}

	log.Println("Starting the server on port", app.Server.Addr)

	err = app.Server.ListenAndServe()
	if err != nil {
		log.Fatalln("couldn't serve on port 4000")
	}
}
