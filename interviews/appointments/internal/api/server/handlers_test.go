package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/ericbutera/appointments/internal/api/server"
	"github.com/ericbutera/appointments/internal/repo"
	"github.com/ericbutera/appointments/internal/test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	AvailabilityUrl = "/availability?trainer_id=%d&starts_at=%s&ends_at=%s"
)

func TestListAvailability(t *testing.T) {
	s := test.NewSetup(t)
	//s.Seed(t)

	gin.SetMode(gin.TestMode)
	handlers, err := server.NewHandlers(s.Repo, s.Location)
	require.NoError(t, err)
	router := server.NewRouter(handlers)

	start := time.Date(2025, 01, 14, repo.BusinessStartHour, 00, 0, 0, s.Location)
	end := start.Add(repo.AppointmentDuration)

	recorder := httptest.NewRecorder()
	url := fmt.Sprintf("/availability?trainer_id=%d&starts_at=%s&ends_at=%s",
		test.TestTrainerID,
		url.QueryEscape(start.Format(server.DateFormat)),
		url.QueryEscape(end.Format(server.DateFormat)),
	)
	request := httptest.NewRequest(http.MethodGet, url, nil)
	router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)

	expected := server.ListAvailabilityResponse{
		Availability: []*server.ListAppointment{
			{
				StartsAt: start,
				EndsAt:   end,
			},
		},
	}
	var actual server.ListAvailabilityResponse
	unmarshal(t, recorder, &actual)

	for _, a := range actual.Availability { // TODO: find a better way to unmarshal time using location
		a.StartsAt = a.StartsAt.In(s.Location)
		a.EndsAt = a.EndsAt.In(s.Location)
	}
	assert.Equal(t, expected, actual)
}

// TODO create table tests to validate list availability inputs
// validations:
// invalid starts_at, ends at
// invalid trainer
// starts at before ends at
// starts at after business hours
// rules:
// no availability

func unmarshal(t *testing.T, recorder *httptest.ResponseRecorder, v any) {
	t.Helper()
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), v))
}

func TestCreateAppointment(t *testing.T) {
	t.Skip("TODO")
}

func TestListAppointment(t *testing.T) {
	t.Skip("TODO")
}
