package server

import (
	"1linkshare/internal/utils"
	"fmt"

	"1linkshare/internal/network"
	"net/http"
	"path/filepath"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (s *Server) HomeHandler(c echo.Context) error {
	return c.Redirect(http.StatusSeeOther, "/upload")
}

func (s *Server) UploadHandler(c echo.Context) error {
	var id string

	if c.Request().Method == http.MethodGet {
		return c.Render(http.StatusOK, "upload.html", nil)
	}

	form, err := c.MultipartForm()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error parsing form data")
	}

	files := form.File["uploadFile"]
	if len(files) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "No file uploaded")
	}

	for _, file := range files {
		id = uuid.New().String()
		filePath := filepath.Join("./uploads", id+"_"+file.Filename)

		src, err := file.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error retrieving the file")
		}
		defer src.Close()

		dst, err := utils.CreateFile(filePath)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error creating the file")
		}
		defer dst.Close()

		if _, err := dst.ReadFrom(src); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error saving the file")
		}

		s.store.Add(id, filePath)
	}

	urlCh := make(chan string)
	go network.RunSSHCommand(urlCh, []string{"-R", "80:localhost:8080", "serveo.net"})
	forwardedURL := <-urlCh

	body := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF8800")).Render
	fmt.Println()
	bigSwag := lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff88")).Padding(0, 1).Render
	urlDisplay := body(fmt.Sprintf("%s/share/%s", forwardedURL, id))
	tunelURLdisplay := "Tunnel URL:  "
	fmt.Println(bigSwag(tunelURLdisplay), body(urlDisplay))
	fmt.Println()
	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("%s/share/%s", forwardedURL, id))
}

func (s *Server) ShareHandler(c echo.Context) error {
	id := c.Param("id")
	filePath, ok := s.store.Get(id)
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "File not found")
	}
	return c.File(filePath)
}

func (s *Server) HealthHandler(c echo.Context) error {
	resp := map[string]interface{}{
		"status":  "OK",
		"version": "1.0.0",
		"time":    time.Now().Format(time.RFC3339),
		"uptime":  time.Since(s.startTime).String(),
	}
	return c.JSON(http.StatusOK, resp)
}
