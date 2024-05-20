package server

import (
	"1linkshare/internal/utils"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

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
		fmt.Println("FILE_ATH " + filePath)

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

	fmt.Printf("test rout : %s", fmt.Sprintf("/share/%s", id))

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/share/%s", id))
}

func (s *Server) ShareHandler(c echo.Context) error {
	id := c.Param("id")
	fmt.Println("TEST: id" + id)
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
