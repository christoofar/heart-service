package floors

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	common "heart-service/common"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", Add)
	return router
}

func Add(w http.ResponseWriter, r *http.Request) {
	var response common.HealthResponse

	decoder := json.NewDecoder(r.Body)
	var data common.HealthData
	err := decoder.Decode(&data)
	if err != nil {
		log.Fatal("Bad floors climbed data was sent")
		response.Message = "Bad floors climbed data was sent"
		response.Time = time.Now()
		render.JSON(w, r, response)
		return
	}

	log.Printf("Floors data received: %s\n", strconv.FormatFloat(data.Readings[0].Value, 'E', -1, 32))

	common.PostValues(data, "Floors Climbed", 0)

	response.Message = "OK"
	response.Time = time.Now()
	render.JSON(w, r, response)
}
