package server

// API times are in the user's timezone (default to PST).
// TODO: make error messages user friendly and consistent

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ericbutera/appointments/internal/repo"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

const (
	DateFormat = "2006-01-02T15:04:05-07:00"
)

type Handlers struct {
	repo     *repo.Repo
	location *time.Location
}

func NewHandlers(r *repo.Repo, location *time.Location) (*Handlers, error) {
	return &Handlers{
		location: location,
		repo:     r,
	}, nil
}

type ErrorResponse struct {
	Error string `example:"unable to fetch feeds" json:"error"`
	// TODO: add error code
	// TODO: add request ID for correlating traces, logs, etc
}

type ListAppointment struct {
	StartsAt  time.Time `json:"starts_at"`
	EndsAt    time.Time `json:"ends_at"`
	UserID    uint      `json:"user_id"`
	TrainerID uint      `json:"trainer_id"`
}

type ListAvailabilityResponse struct {
	Availability []*ListAppointment `json:"availability"`
}

// TODO requires user auth
// TODO: paginate
func (h *Handlers) ListAvailability(c *gin.Context) {
	trainerID, err := strconv.Atoi(c.Query("trainer_id"))
	if err != nil || trainerID <= 0 { // TODO: max bound
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid trainer_id"})
		return
	}

	start, err := stringToTime(c.Query("starts_at"), h.location)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid starts_at date"})
		return
	}

	end, err := stringToTime(c.Query("ends_at"), h.location)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid ends_at date"})
		return
	}

	if start.After(end) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "starts_at time must be before ends_at"})
		return
	}

	if end.Sub(start) > repo.MaxAvailabilityRange {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: repo.ErrExceedsMaxRange.Error()})
		return
	}

	// UX improvements:
	// TODO: if invalid start, auto change to current time
	// TODO: if start == end auto change to +2 weeks

	res, err := h.repo.GetAvailability(c.Request.Context(), uint(trainerID), start, end)
	if err != nil {
		if errors.Is(err, repo.ErrNoAvailability) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "no availability"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "unable to get availability"})
		return
	}

	c.JSON(http.StatusOK, ListAvailabilityResponse{
		Availability: toListAppointment(res),
	})
}

type CreateAppointment struct {
	StartsAt  time.Time `binding:"required" json:"starts_at"`
	EndsAt    time.Time `binding:"required" json:"ends_at"`
	UserID    uint      `binding:"required,gt=0" json:"user_id"` // TODO set via auth
	TrainerID uint      `binding:"required,gt=0" json:"trainer_id"`
}

type CreateAppointmentRequest struct {
	Appointment *CreateAppointment `json:"appointment" binding:"required"`
}

type CreateAppointmentResponse struct {
	ID uint `json:"id"`
}

// TODO: prevent creating new appointments in the past
// TODO requires user auth
func (h *Handlers) CreateAppointment(c *gin.Context) {
	var req CreateAppointmentRequest
	if err := c.ShouldBind(&req); err != nil { // Note: basic validation, domain validation is in the repo
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: fmt.Sprintf("validation error: %s", err.Error()),
		})
		return
	}

	appointment := &repo.Appointment{
		StartsAt:  req.Appointment.StartsAt,
		EndsAt:    req.Appointment.EndsAt,
		UserID:    req.Appointment.UserID,    // TODO use auth user id
		TrainerID: req.Appointment.TrainerID, // TODO verify user is allowed to book this trainer
	}
	res, err := h.repo.CreateAppointment(c.Request.Context(), appointment)
	if err != nil {
		if errors.As(err, &repo.ValidationErrors{}) {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid appointment: " + err.Error()})
			return
		} else if err == repo.ErrDuplicateAppointment {
			c.JSON(http.StatusConflict, ErrorResponse{Error: "time unavailable"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "unable to create appointment"})
		return
	}

	c.JSON(http.StatusOK, CreateAppointmentResponse{
		ID: res.ID,
	})
}

type ListAppointmentsResponse struct {
	Appointments []*ListAppointment `json:"appointments"`
}

// TODO requires trainer auth
// TODO: paginate
func (h *Handlers) ListAppointments(c *gin.Context) {
	trainerID, err := strconv.Atoi(c.Query("trainer_id")) // TODO use auth trainer id
	if err != nil || trainerID <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid trainer_id"})
		return
	}

	start, err := stringToTime(c.Query("starts_at"), h.location)
	if err != nil {
		start = time.Now().In(h.location).Truncate(time.Hour)
	}
	end := time.Now().In(h.location).Add(2 * repo.Week) // this is an assumption

	res, err := h.repo.GetAppointments(c.Request.Context(), uint(trainerID), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "unable to get appointments"})
		return
	}

	data := toListAppointment(res)

	c.JSON(http.StatusOK, ListAppointmentsResponse{
		Appointments: data,
	})
}

func toListAppointment(res []*repo.Appointment) []*ListAppointment {
	// TODO: use a mapping layer for transforms (goverter)
	availability := make([]*ListAppointment, 0, len(res))
	for _, a := range res {
		availability = append(availability, &ListAppointment{
			StartsAt:  a.StartsAt,
			EndsAt:    a.EndsAt,
			UserID:    a.UserID,
			TrainerID: a.TrainerID,
		})
	}
	return availability
}

type WorkoutStatus struct {
	ScheduledAt time.Time
	Status      int // TODO: try enum
}

type WorkoutStatusResult struct {
	Data           []WorkoutStatus
	Conclusion     string // TODO: enum
	SuccessPercent float32
}

const WorkoutApi = "https://example.com"
const RequiredPercentage float32 = 80

type WorkoutApiWorkout struct {
	Type            string    `json:"type"` // TODO: use oneof validator
	ScheduledAt     time.Time `json:"scheduled_at"`
	CompletionState string    `json:"completion_state"`
}

func (h *Handlers) WorkoutStatus(c *gin.Context) {
	// userID string /* TODO filter for time range*/
	// assume we are a trainer trainerID :=
	userID := c.Param("userID") // TODO: validate UUID

	var workouts []WorkoutApiWorkout
	client := resty.New()
	resp, err := client.R().
		SetResult(&workouts).
		Get(fmt.Sprintf(WorkoutApi, userID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "unable to fetch feeds"})
		return // TODO: dont show raw error to user
	}

	if resp.StatusCode() != http.StatusOK {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "unable to fetch feeds"})
		return
	}

	// TODO: do not use magic strings for values like rest, workout, completed etc
	// TODO: sort
	workoutResults := make([]WorkoutStatus, 0, len(workouts))

	successCounter := 0
	for _, workout := range workouts { // TODO: make mapper to change type
		status := GetStatus(workout)
		if status == 1 {
			successCounter++
		}

		workoutResults = append(workoutResults, WorkoutStatus{
			ScheduledAt: workout.ScheduledAt,
			Status:      status,
		})
	}

	// make a threshold for deciding good or problem
	// if 80% of the workouts are completed then good else problem
	// successCounter * 100 / count = percentage

	percentage := CalculatePercentage(successCounter, len(workouts))

	c.JSON(http.StatusOK, WorkoutStatusResult{
		Data:           workoutResults,
		Conclusion:     string(GetConclusions(percentage)),
		SuccessPercent: percentage,
	})
}

func CalculatePercentage(successCounter int, total int) float32 {
	return (float32(successCounter) * 100) / float32(total)
}

type WorkoutConclusion string

const (
	WorkoutConclusionSuccess WorkoutConclusion = "success"
	WorkoutConclusionProblem WorkoutConclusion = "problem"
)

func GetConclusions(percentage float32) WorkoutConclusion {
	isAboveThreshold := percentage >= RequiredPercentage
	if isAboveThreshold {
		return WorkoutConclusionSuccess
	}
	return WorkoutConclusionProblem
}

type WorkoutTypes string

const (
	WorkoutTypesRest    WorkoutTypes = "rest"
	WorkoutTypesWorkout WorkoutTypes = "workout"
)

type WorkoutStatuses string

const (
	WorkoutStatusCompleted WorkoutStatuses = "completed"
	WorkoutStatusPartial   WorkoutStatuses = "partial"
	WorkoutStatusNone      WorkoutStatuses = "none"
)

func GetStatus(workout WorkoutApiWorkout) int {
	status := 0
	workoutType := WorkoutTypes(workout.Type)
	workoutCompletionState := WorkoutStatuses(workout.CompletionState)
	if workoutType == WorkoutTypesRest { // TODO: extract function (should be in library or business logic as this will change over time)
		status = 1
	} else if workoutType == WorkoutTypesWorkout {
		if workoutCompletionState == WorkoutStatusCompleted {
			status = 1
		} else if workoutCompletionState == WorkoutStatusPartial {
			status = 1
		}
	}
	return status
}

func stringToTime(val string, location *time.Location) (time.Time, error) {
	parsed, err := time.ParseInLocation(DateFormat, val, location)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid datetime format: %v", err)
	}
	parsed = parsed.Truncate(time.Minute) // remove seconds, nanoseconds
	return parsed, nil
}
