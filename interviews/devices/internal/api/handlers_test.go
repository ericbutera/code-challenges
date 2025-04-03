package api_test

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/ericbutera/project/internal/api"
	"github.com/ericbutera/project/internal/db"
	"github.com/ericbutera/project/internal/models"
	"github.com/ericbutera/project/internal/repo"
	"github.com/ericbutera/project/internal/test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	DevicesReadingEndpoint      = "/v1/devices/%s/readings"
	DevicesReadingCountEndpoint = "/v1/devices/%s/readings/count"
	DevicesReadingsLatest       = "/v1/devices/%s/readings/latest"
)

type testSetup struct {
	handlers *api.Handlers
	repo     *repo.Repo
	router   *gin.Engine
	db       *db.MockDB
}

func setup(t *testing.T) *testSetup {
	t.Helper()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	slog.SetDefault(logger)

	gin.SetMode(gin.TestMode)

	db := new(db.MockDB)

	r, err := repo.New(db)
	require.NoError(t, err)

	handlers, err := api.NewHandlers(r)
	require.NoError(t, err)

	router := api.NewRouter(handlers)

	return &testSetup{
		handlers: handlers,
		repo:     r,
		router:   router,
		db:       db,
	}
}

func TestDeviceReadingCount(t *testing.T) {
	s := setup(t)

	s.db.EXPECT().
		GetReadingCountByDevice(test.TestDeviceID).
		Return(&db.DeviceReadingsCount{Count: 17}, nil)

	url := fmt.Sprintf(DevicesReadingCountEndpoint, test.TestDeviceID)
	request := httptest.NewRequest(http.MethodGet, url, nil)
	recorder := httptest.NewRecorder()
	s.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)

	var data map[string]any
	decoder := json.NewDecoder(recorder.Body)
	decoder.UseNumber()
	require.NoError(t, decoder.Decode(&data))
	assert.Equal(t, json.Number("17"), data["cumulative_count"])
}

// TODO: test reading count with no data

func TestDeviceLatestReading(t *testing.T) {
	s := setup(t)

	dt := "2021-09-01T17:00:00-05:00"
	val, err := api.TimeFromString(dt)
	require.NoError(t, err)

	s.db.EXPECT().
		GetLatestReadingByDevice(test.TestDeviceID).
		Return(&db.DeviceLatestReading{
			Timestamp: val,
		}, nil)

	url := fmt.Sprintf(DevicesReadingsLatest, test.TestDeviceID)
	request := httptest.NewRequest(http.MethodGet, url, nil)
	recorder := httptest.NewRecorder()
	s.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)

	var data map[string]any
	unmarshal(t, recorder, &data)
	assert.Equal(t, dt, data["latest_timestamp"])
}

// TODO: test latest reading with no data

func TestDeviceStoreReadings(t *testing.T) {
	s := setup(t)

	raw := `
	{
		"id": "36d5658a-6908-479e-887e-a949ec199272",
		"readings": [
			{"timestamp":"2021-09-01T17:00:00-05:00","count":17}
		]
	}
	`
	ts, err := api.TimeFromString("2021-09-01T17:00:00-05:00")
	require.NoError(t, err)

	s.db.EXPECT().
		StoreDeviceReadings(test.TestDeviceID, []*models.Reading{
			{Timestamp: ts, Count: 17},
		}).
		Return(&db.StoreDeviceReadingsResult{}, nil)

	url := fmt.Sprintf(DevicesReadingEndpoint, test.TestDeviceID)
	request := httptest.NewRequest(http.MethodPost, url, strings.NewReader(raw))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	s.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestDeviceStoreReadings_Invalid(t *testing.T) {
	s := setup(t)

	cases := []struct {
		name string
		err  string
		raw  string
	}{
		{"invalid device id", "id", `{"id": 12345, "readings": []}`},
		{"invalid timestamp", "cannot parse", `{"id": "36d5658a-6908-479e-887e-a949ec199272", "readings": [{"timestamp":"asdf-09-01T17:00:00-05:00","count":"hello"}]}`},
		{"invalid count", "of type int", `{"id": "36d5658a-6908-479e-887e-a949ec199272", "readings": [{"timestamp":"2021-09-01T17:00:00-05:00","count":"hello"}]}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf(DevicesReadingEndpoint, test.TestDeviceID)
			request := httptest.NewRequest(http.MethodPost, url, strings.NewReader(tc.raw))
			request.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			s.router.ServeHTTP(recorder, request)
			require.Equal(t, http.StatusBadRequest, recorder.Code)
		})
	}
}

func TestDeviceStore_InvalidContent(t *testing.T) {
	s := setup(t)
	url := fmt.Sprintf(DevicesReadingEndpoint, test.TestDeviceID)
	request := httptest.NewRequest(http.MethodPost, url, strings.NewReader("invalid content"))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	s.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestDeviceStore_InvalidContentType(t *testing.T) {
	s := setup(t)
	url := fmt.Sprintf(DevicesReadingEndpoint, test.TestDeviceID)
	request := httptest.NewRequest(http.MethodPost, url, strings.NewReader("invalid content"))
	request.Header.Set("Content-Type", "text/plain")
	recorder := httptest.NewRecorder()
	s.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func unmarshal(t *testing.T, recorder *httptest.ResponseRecorder, v any) {
	t.Helper()
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), v))
}
