package handlers

import (
	"net/http"

	"github.com/jcsw/go-api-learn/pkg/infra/database"
)

type monitorComponent struct {
	Component string `json:"component"`
	Status    string `json:"status"`
}

// MonitorHandler function to handle "/monitor"
func MonitorHandler(w http.ResponseWriter, r *http.Request) {
	monitors := []monitorComponent{}
	monitors = append(monitors, retriveMongoDBStatus())
	respondWithJSON(w, http.StatusOK, monitors)
}

func retriveMongoDBStatus() monitorComponent {
	mongoDBStatus := monitorComponent{Component: "MongoDB"}
	if database.IsMongoClientAlive() {
		mongoDBStatus.Status = "OK"
	} else {
		mongoDBStatus.Status = "ERROR"
	}

	return mongoDBStatus
}
