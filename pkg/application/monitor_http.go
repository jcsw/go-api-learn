package application

import (
	"net/http"
	"sync/atomic"

	"github.com/jcsw/go-api-learn/pkg/infra/database"
)

type monitorComponent struct {
	Component string `json:"component"`
	Status    string `json:"status"`
}

// MonitorHandle function to handle "/monitor"
func MonitorHandle(w http.ResponseWriter, r *http.Request) {

	monitors := []monitorComponent{}

	if r.Method == "GET" {
		monitors = append(monitors, mongoDBStatus())

		respondWithJSON(w, http.StatusOK, monitors)
		return
	}

	respondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
}

func mongoDBStatus() monitorComponent {

	mongoDBStatus := monitorComponent{Component: "MongoDB"}

	if atomic.LoadInt32(database.GetMongoDBStatus()) == 1 {
		mongoDBStatus.Status = "OK"
	} else {
		mongoDBStatus.Status = "ERROR"
	}

	return mongoDBStatus
}
