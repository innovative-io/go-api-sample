package controllers

import (
	"net/http"
)

// @Summary Gets the status of the server
// @Description gets the status of the server
// @Produce  json
// @Success 200 {object} map[string]string "ok"
// @Failure 503 {string} string "database unavailable"
// @Router /health [get]
func (r *Router) HealthGet(w http.ResponseWriter, req *http.Request) {
	sqlDB, err := r.db.DB()
	if err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"message": "database unavailable"})
		return
	}
	if err := sqlDB.PingContext(req.Context()); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"message": "database unavailable"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "ok"})
}
