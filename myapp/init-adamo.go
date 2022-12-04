package main

import (
	"log"
	"myapp/handlers"
	"os"

	"github.com/aag2807/adamo-framework"
)

func initApp() *application {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// initialize adamos
	adam := &adamo.Adamo{}
	err = adam.New(path)
	if err != nil {
		log.Fatal(err)
	}

	adam.AppName = "myapp"

	myHandlers := &handlers.Handlers{
		App: adam,
	}

	app := &application{
		App:      adam,
		Handlers: myHandlers,
	}

	app.App.Routes = app.routes()

	return app
}
