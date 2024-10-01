package handlers

import (
	"fmt"
	"log"
	"net/http"
)

type App struct {
	l *log.Logger
}

func NewApp(l *log.Logger) *App {
	return &App{l}
}

func (app *App) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	app.l.Println("Request received")
	fmt.Fprint(w, "Hello World")
}
