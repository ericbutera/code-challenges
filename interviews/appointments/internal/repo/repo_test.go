package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/ericbutera/appointments/internal/repo"
	"github.com/ericbutera/appointments/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAvailability(t *testing.T) {
	s := test.NewSetup(t)

	start := time.Date(2025, 01, 13, repo.BusinessStartHour, 0, 0, 0, s.Location)
	end := time.Date(2025, 01, 13, repo.BusinessEndHour, 0, 0, 0, s.Location)

	appointments, err := s.Repo.GetAvailability(context.Background(), test.TestTrainerID, start, end)
	require.NoError(t, err)

	require.Len(t, appointments, 18)
	assert.Equal(t, appointments[0].StartsAt, start)
	assert.Equal(t, appointments[17].EndsAt, end)
}

func TestGetAvailability_WithExistingAppointment(t *testing.T) {
	s := test.NewSetup(t)

	start := time.Date(2025, 01, 13, repo.BusinessStartHour, 0, 0, 0, s.Location)
	end := start.Add(repo.AppointmentDuration * 2)

	availability, err := s.Repo.GetAvailability(context.Background(), test.TestTrainerID, start, end)
	require.NoError(t, err)
	require.Len(t, availability, 2)

	slot := availability[0]
	appointment := repo.Appointment{
		StartsAt:  slot.StartsAt,
		EndsAt:    slot.EndsAt,
		UserID:    test.TestUserID,
		TrainerID: test.TestTrainerID,
	}
	require.NoError(t, s.DB.Create(&appointment).Error)

	availability2, err := s.Repo.GetAvailability(context.Background(), test.TestTrainerID, start, end)
	require.NoError(t, err)
	require.Len(t, availability2, 1)
}

// TODO: test GetAvailability supports multiple existing appointments

func TestIsAvailable(t *testing.T) {
	s := test.NewSetup(t)

	start := time.Date(2025, 01, 13, repo.BusinessStartHour, 0, 0, 0, s.Location)
	appointment := repo.Appointment{
		StartsAt:  start,
		EndsAt:    start.Add(repo.AppointmentDuration),
		UserID:    test.TestUserID,
		TrainerID: test.TestTrainerID,
	}

	ok, err := s.Repo.IsAvailable(context.Background(), test.TestTrainerID, appointment.StartsAt, appointment.EndsAt)
	require.NoError(t, err)
	assert.True(t, ok)

	require.NoError(t, s.DB.Create(&appointment).Error)

	ok, err = s.Repo.IsAvailable(context.Background(), test.TestTrainerID, appointment.StartsAt, appointment.EndsAt)
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestCreateAppointment(t *testing.T) {
	s := test.NewSetup(t)

	start := time.Date(2025, 01, 13, repo.BusinessStartHour, 0, 0, 0, s.Location)
	appointment := repo.Appointment{
		StartsAt:  start,
		EndsAt:    start.Add(repo.AppointmentDuration),
		UserID:    test.TestUserID,
		TrainerID: test.TestTrainerID,
	}

	res, err := s.Repo.CreateAppointment(context.Background(), &appointment)
	require.NoError(t, err)

	var actual repo.Appointment
	require.NoError(t, s.DB.First(&actual, res.ID).Error)
	test.Diff(t, appointment, actual)
}

func TestCreateAppointment_InvalidTime(t *testing.T) {
	s := test.NewSetup(t)

	cases := []struct {
		name     string
		time     time.Time
		duration time.Duration
		isErr    bool
	}{
		{
			name:     "start time is 00",
			time:     time.Date(2025, 01, 13, repo.BusinessStartHour, 0, 0, 0, s.Location),
			duration: repo.AppointmentDuration,
			isErr:    false,
		},
		{
			name:     "start time is 30",
			time:     time.Date(2025, 01, 13, repo.BusinessStartHour, 30, 0, 0, s.Location),
			duration: repo.AppointmentDuration,
			isErr:    false,
		},
		{
			name:     "end time is 00",
			time:     time.Date(2025, 01, 13, repo.BusinessEndHour-1, 0, 00, 0, s.Location),
			duration: repo.AppointmentDuration,
			isErr:    false,
		},
		{
			name:     "end time is 30",
			time:     time.Date(2025, 01, 13, repo.BusinessEndHour-1, 30, 00, 0, s.Location),
			duration: repo.AppointmentDuration,
			isErr:    false,
		},
		{
			name:     "start time is not on the hour or half hour",
			time:     time.Date(2025, 01, 13, repo.BusinessStartHour, 01, 0, 0, s.Location),
			duration: repo.AppointmentDuration,
			isErr:    true,
		},
		{
			name:     "invalid duration",
			time:     time.Date(2025, 01, 13, repo.BusinessStartHour, 01, 0, 0, s.Location),
			duration: 1 * time.Hour,
			isErr:    true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			appointment := repo.Appointment{
				StartsAt:  tc.time,
				EndsAt:    tc.time.Add(tc.duration),
				UserID:    test.TestUserID,
				TrainerID: test.TestTrainerID,
			}
			_, err := s.Repo.CreateAppointment(context.Background(), &appointment)
			if tc.isErr {
				require.Error(t, err, tc.name)
			} else {
				require.NoError(t, err, tc.name)
			}
		})
	}
}

func TestAppointmentExists(t *testing.T) {
	s := test.NewSetup(t)

	start := time.Date(2025, 01, 13, repo.BusinessStartHour, 0, 0, 0, s.Location)
	appointment := repo.Appointment{
		StartsAt:  start,
		EndsAt:    start.Add(repo.AppointmentDuration),
		UserID:    test.TestUserID,
		TrainerID: test.TestTrainerID,
	}
	require.NoError(t, s.DB.Create(&appointment).Error)

	appointment2 := appointment
	appointment2.ID = 0

	_, err := s.Repo.CreateAppointment(context.Background(), &appointment2)
	require.Error(t, err)
	assert.ErrorIs(t, err, repo.ErrDuplicateAppointment)
}

func TestFindOpening(t *testing.T) {
	s := test.NewSetup(t)

	cases := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "first appointment of the day",
			input:    time.Date(2025, 01, 13, repo.BusinessStartHour, 0, 0, 0, s.Location), // 8:00 AM
			expected: time.Date(2025, 01, 13, repo.BusinessStartHour, 0, 0, 0, s.Location), // 8:00 AM
		},
		{
			name:     "last appointment of the day",
			input:    time.Date(2025, 01, 13, repo.BusinessEndHour-1, repo.AppointmentMinutes, 0, 0, s.Location), // 4:30 PM
			expected: time.Date(2025, 01, 13, repo.BusinessEndHour-1, repo.AppointmentMinutes, 0, 0, s.Location), // 4:30 PM
		},
		{
			name:     "time before open moves to first appointment of the day",
			input:    time.Date(2025, 01, 13, repo.BusinessStartHour-1, 0, 0, 0, s.Location), // 7:00 AM
			expected: time.Date(2025, 01, 13, repo.BusinessStartHour, 0, 0, 0, s.Location),   // 8:00 AM
		},
		{
			name:     "second appointment of the day",
			input:    time.Date(2025, 01, 13, repo.BusinessStartHour, repo.AppointmentMinutes, 0, 0, s.Location), // 8:30 AM
			expected: time.Date(2025, 01, 13, repo.BusinessStartHour, repo.AppointmentMinutes, 0, 0, s.Location), // 8:30 AM
		},
		{
			name:     "after hours moves to next business day",
			input:    time.Date(2025, 01, 13, repo.BusinessEndHour, 0, 0, 0, s.Location),   // Monday 5:00 PM
			expected: time.Date(2025, 01, 14, repo.BusinessStartHour, 0, 0, 0, s.Location), // Tuesday 8:00 AM
		},
		{
			name:     "Saturday moves to Monday",
			input:    time.Date(2025, 01, 11, repo.BusinessStartHour, 0, 0, 0, s.Location), // Saturday 8:00 AM
			expected: time.Date(2025, 01, 13, repo.BusinessStartHour, 0, 0, 0, s.Location), // Monday 8:00 AM
		},
		{
			name:     "Sunday moves to Monday",
			input:    time.Date(2025, 01, 12, repo.BusinessStartHour, 0, 0, 0, s.Location), // Sunday 8:00 AM
			expected: time.Date(2025, 01, 13, repo.BusinessStartHour, 0, 0, 0, s.Location), // Monday 8:00 AM
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			start, err := s.Repo.FindOpening(tc.input, tc.input.Add(repo.Week))
			require.NoError(t, err)
			assert.Equal(t, tc.expected, start, tc.name)
		})
	}
}
