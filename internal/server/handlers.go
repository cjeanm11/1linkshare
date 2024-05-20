package server

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)


func (s *Server) HealthHandler(c echo.Context) error {
	resp := map[string]interface{}{
		"status":  "OK",
		"version": "1.0.0",
		"time":    time.Now().Format(time.RFC3339),
		"uptime":  time.Since(s.startTime).String(),
	}
	return c.JSON(http.StatusOK, resp)
}