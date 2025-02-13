package monitoring

import (
	"fmt"
	"net/http"
	"time"

	"github.com/turgaysozen/algotrading/db"
	"github.com/turgaysozen/algotrading/redisclient"
	"github.com/turgaysozen/algotrading/wsclient"
)

func ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	err := db.Database.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status": "not ready", "reason": "database unreachable"}`)
		return
	}

	err = redisclient.RedisHealth()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status": "not ready", "reason": "redis unreachable"}`)
		return
	}

	if !wsclient.Connected {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status": "not ready", "reason": "WebSocket unreachable"}`)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := HealthCheckResponse{
		Status: "ready",
		Time:   time.Now().String(),
	}
	fmt.Fprintf(w, `{"status": "%s", "time": "%s"}`, response.Status, response.Time)
}
