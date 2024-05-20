package tests

import (
	"net/http"
	"net/http/httptest"
	"1linkshare/internal/server"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	resp := httptest.NewRecorder()
	c := e.NewContext(req, resp)
	s := &server.Server{}
	if err := s.HealthHandler(c); err != nil {
		t.Errorf("handler() error = %v", err)
		return
	}
	if resp.Code != http.StatusOK {
		t.Errorf("handler() wrong status code = %v", resp.Code)
		return
	}
	if assert.NoError(t, s.HealthHandler(c)) {
		assert.Equal(t, http.StatusOK, resp.Code)
	}
}
