package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ericbutera/project/internal/models"
	"github.com/ericbutera/project/internal/repo"
	"github.com/gin-gonic/gin"
)

const DateFormatISO8601 = "2006-01-02T15:04:05-07:00"

type Handlers struct {
	repo *repo.Repo
}

func NewHandlers(r *repo.Repo) (*Handlers, error) {
	return &Handlers{
		repo: r,
	}, nil
}

type ErrorResponse struct {
	Error string `json:"error"`
	// TODO: add error code
	// TODO: add request ID for correlating traces, logs, etc
}

type DeviceStoreReadingsRequest struct {
	DeviceID string            `binding:"required,uuid4" description:"Device ID" json:"id"`
	Readings []*models.Reading `binding:"required,dive"  description:"Readings"  json:"readings"` // TODO: dive validation
}

type DeviceStoreReadingsResponse struct {
	CreateCount int `description:"number of readings that were created" json:"create_count"`
}

// store readings for a device; will create a device if it doesn't exist.
func (h *Handlers) DeviceStoreReadings(c *gin.Context) {
	// TODO: validate path device matches body device
	// TODO: go validation library supports "friendly" error messages (i11n becomes a concern)
	// TODO: throw out batch if threshold of invalid readings are reached
	var req DeviceStoreReadingsRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: fmt.Sprintf("validation error: %s", err.Error()),
		})
		return
	}

	_, err := h.repo.StoreReadings(c.Request.Context(), req.DeviceID, req.Readings)
	if err != nil {
		// TODO: support validation errors; if errors.As(err, &repo.ValidationErrors{}) { c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid record: " + err.Error()}) return }
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "unable to save readings"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

type DeviceLatestReadingResponse struct {
	LatestTimestamp string `json:"latest_timestamp"`
}

func (h *Handlers) DeviceLatestReading(c *gin.Context) {
	res, err := h.repo.GetLatestReadingByDevice(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "device not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "unable to get cumulative count"})
		return
	}
	c.JSON(http.StatusOK, DeviceLatestReadingResponse{
		LatestTimestamp: TimeToString(res.LatestReading),
	})
}

type DeviceReadingsCountResponse struct {
	Count int64 `json:"cumulative_count"`
}

func (h *Handlers) DeviceReadingCount(c *gin.Context) {
	// TODO: 404 if device not found
	res, err := h.repo.GetCountByDevice(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "unable to get latest reading"})
		return
	}
	c.JSON(http.StatusOK, DeviceReadingsCountResponse{
		Count: res.Count,
	})
}

type DeviceGetReadingsResponse struct {
	Readings JSONSlice[models.Reading] `json:"readings"`
}

func (h *Handlers) DeviceGetReadings(c *gin.Context) {
	res, err := h.repo.GetReadingsByDevice(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "unable to get readings"})
		return
	}
	c.JSON(http.StatusOK, &DeviceGetReadingsResponse{
		Readings: JSONSlice[models.Reading](res),
	})
}

func TimeFromString(s string) (time.Time, error) {
	return time.Parse(DateFormatISO8601, s)
}

func TimeToString(t time.Time) string {
	return t.Format(DateFormatISO8601)
}

// Prevents empty slices from being marshaled as `null` which can break some clients.
type JSONSlice[T any] []T

func (s JSONSlice[T]) MarshalJSON() ([]byte, error) {
	if s == nil {
		return json.Marshal([]T{})
	}
	return json.Marshal([]T(s))
}
