package main

import (
	"myapp/handlers"

	"github.com/aag2807/adamo-framework"
)

type application struct {
	App      *adamo.Adamo
	Handlers *handlers.Handlers
}

func main() {
	a := initApp()
	a.App.ListenAndServe()
}
