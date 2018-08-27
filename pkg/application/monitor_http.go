package application

import (
	"net/http"

	"github.com/jcsw/go-api-learn/pkg/infra/database"
	"gopkg.in/macaron.v1"
)

type monitorComponent struct {
	Component string `json:"component"`
	Status    string `json:"status"`
}

// MonitorHandle function to handle "/monitor"
func MonitorHandle(ctx *macaron.Context) {
	monitors := []monitorComponent{}
	monitors = append(monitors, getMongoDBStatus())
	respondWithJSON(ctx, http.StatusOK, monitors)
}

func getMongoDBStatus() monitorComponent {
	mongoDBStatus := monitorComponent{Component: "MongoDB"}
	if database.IsMongoDBSessionAlive() {
		mongoDBStatus.Status = "OK"
	} else {
		mongoDBStatus.Status = "ERROR"
	}

	return mongoDBStatus
}
