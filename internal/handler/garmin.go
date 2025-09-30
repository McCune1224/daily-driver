package handler

import (
	"daily-driver/internal/db"
	"daily-driver/web/static/templates"
	"io"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (h *Handler) RenderGarminUploadPage(c echo.Context) error {
	pageStr := c.QueryParam("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	const limit = 10
	offset := (page - 1) * limit

	q := db.New(h.DBPool)
	files, err := q.ListGarminFilesPaginated(c.Request().Context(), db.ListGarminFilesPaginatedParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		h.Logger.Error("Failed to get garmin fit files", zap.Error(err))
		return c.String(500, "Failed to get garmin fit files")
	}

	total, err := q.CountGarminFiles(c.Request().Context())
	if err != nil {
		h.Logger.Error("Failed to count garmin fit files", zap.Error(err))
		return c.String(500, "Failed to count garmin fit files")
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	return Render(c, 200, templates.Admin(files, page, int(totalPages), total))
}

func (h *Handler) RenderGarminPanel(c echo.Context) error {
	return Render(c, 200, templates.PanelGarmin(1))
}

func (h *Handler) UploadGarminFile(c echo.Context) error {
	// grab the file from the multipart/form-data

	form, err := c.MultipartForm()
	if err != nil {
		h.Logger.Error("Failed to get file from form", zap.Error(err))
		return c.String(400, "Failed to get file from form")
	}

	files := form.File["file"]

	for _, file := range files {
		// make the file aviabale to be opened / read:
		src, err := file.Open()
		if err != nil {
			h.Logger.Error("Failed to open uploaded file", zap.Error(err))
			return c.String(500, "Failed to open uploaded file")
		}
		defer src.Close()

		//get the byte data from the file
		data, err := io.ReadAll(src)
		if err != nil {
			h.Logger.Error("Failed to read uploaded file", zap.Error(err))
			return c.String(500, "Failed to read uploaded file")
		}

		q := db.New(h.DBPool)
		fitFile, err := q.InsertGarminFitFile(c.Request().Context(), db.InsertGarminFitFileParams{
			Filename:     file.Filename,
			Data:         data,
			FileCategory: db.NullGarminFileCategory{Valid: true, GarminFileCategory: db.GarminFileCategoryActivity},
		})

		if err != nil {
			h.Logger.Error("Failed to insert garmin fit file", zap.Error(err))
			if strings.Contains(err.Error(), "duplicate key value") {
				h.Logger.Info("Duplicate file upload attempted", zap.String("filename", file.Filename))
				continue
			}
			return c.String(500, "Failed to insert garmin fit file")
		}
		h.Logger.Info("Received file", zap.String("filename", fitFile.Filename), zap.Int64("size", file.Size))
	}

	return c.String(200, "Files uploaded successfully")

}
