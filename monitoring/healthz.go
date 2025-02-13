package monitoring

import (
	"fmt"
	"net/http"
	"time"
)

type HealthCheckResponse struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

func LivenessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := HealthCheckResponse{
		Status: "alive",
		Time:   time.Now().String(),
	}
	fmt.Fprintf(w, `{"status": "%s", "time": "%s"}`, response.Status, response.Time)
}
