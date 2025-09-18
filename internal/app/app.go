package app

import "net/http"

type App struct {
}

func (a *App) HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
