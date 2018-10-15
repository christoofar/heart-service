package ping

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type PingMessage struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", Ping)
	return router
}

func Ping(w http.ResponseWriter, r *http.Request) {
	var ping PingMessage

	ping.Message = "Ping OK"
	ping.Time = time.Now()
	render.JSON(w, r, ping)
}
