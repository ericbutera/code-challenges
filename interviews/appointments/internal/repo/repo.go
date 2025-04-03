package repo

// Database times are stored in UTC.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

const (
	BusinessLocation     = "America/Los_Angeles"
	Day                  = 24 * time.Hour
	Week                 = 7 * Day
	MaxAvailabilityRange = 40 * Day
	AppointmentDuration  = AppointmentMinutes * time.Minute
	AppointmentMinutes   = 30
	BusinessStartHour    = 8  // 8:00 AM PST
	BusinessEndHour      = 17 // 5:00 PM PST (24-hour format)
)

var (
	ErrInvalidDateFormat    = errors.New("date format must be ISO8601")
	ErrExceedsMaxRange      = fmt.Errorf("range between start and end cannot exceed %s", MaxAvailabilityRange) // TODO make user friendly
	ErrNoAvailability       = errors.New("no availability found")
	ErrDuplicateAppointment = errors.New("appointment already exists")
	ErrValidation           = errors.New("validation error")
)

type CreateAppointmentResult struct {
	ID uint `json:"id"`
}

// TODO: separate repo from gorm implementation (allows DI in test/api)
// type Repo interface { ... }
// type GormRepo struct { ... }

type Repo struct {
	db       *gorm.DB
	location *time.Location
}

func New(db *gorm.DB, location *time.Location) (*Repo, error) {
	return &Repo{
		db:       db,
		location: location,
	}, nil
}

type ValidationErrors struct {
	Errors []error `json:"errors"`
}

func (e ValidationErrors) Error() string {
	var msgs []string
	for _, err := range e.Errors {
		msgs = append(msgs, err.Error())
	}
	return fmt.Sprintf("validation failed: %s", strings.Join(msgs, "; "))
}

// Ensures if appointment data is valid, it will reject invalid data, not sanitize it.
func (r *Repo) ValidateAppointment(ctx context.Context, appointment *Appointment) error {
	// TODO: use a validation lib (like go-playground/validator)
	// TODO: prevent creating appointments in the past (but not for this demo because it breaks all the tests)
	var errs []error

	if appointment.StartsAt.Truncate(time.Minute) != appointment.StartsAt {
		errs = append(errs, errors.New("start time must be on the minute"))
	}
	if appointment.EndsAt.Truncate(time.Minute) != appointment.EndsAt {
		errs = append(errs, errors.New("end time must be on the minute"))
	}
	if appointment.StartsAt.After(appointment.EndsAt) {
		errs = append(errs, errors.New("start time must be before end time"))
	}
	if appointment.StartsAt.Minute() != 0 && appointment.StartsAt.Minute() != 30 {
		errs = append(errs, errors.New("start time must be on the hour or half hour"))
	}
	if appointment.EndsAt.Minute() != 0 && appointment.EndsAt.Minute() != 30 {
		errs = append(errs, errors.New("end time must be on the hour or half hour"))
	}
	if appointment.StartsAt == appointment.EndsAt {
		errs = append(errs, errors.New("start time and end time cannot be the same"))
	}
	if appointment.EndsAt.Sub(appointment.StartsAt) != AppointmentDuration {
		errs = append(errs, fmt.Errorf("appointment must be %d minutes", AppointmentMinutes))
	}
	// TODO: these have problems where end fails because appointment.end is actually 5pm
	// separate start and end validations?
	// start allow is 8:00am - 4:30pm
	// end allow is 8:30am - 5:00pm
	if !r.IsBusinessOpen(appointment.StartsAt) {
		errs = append(errs, errors.New("start time must be during business hours"))
	}
	endHack := appointment.EndsAt.Add(-AppointmentDuration)
	if !r.IsBusinessOpen(endHack) {
		errs = append(errs, errors.New("end time must be during business hours"))
	}

	if len(errs) > 0 {
		return ValidationErrors{Errors: errs}
	}

	return nil
}

func (r *Repo) CreateAppointment(ctx context.Context, appointment *Appointment) (*CreateAppointmentResult, error) {
	// New appointments only, no update supported
	if err := r.ValidateAppointment(ctx, appointment); err != nil {
		return nil, err
	}

	var result *CreateAppointmentResult
	err := r.db.Transaction(func(tx *gorm.DB) error {
		available, err := r.IsAvailable(ctx, appointment.TrainerID, appointment.StartsAt, appointment.EndsAt)
		if err != nil {
			return err
		}
		if !available {
			return ErrDuplicateAppointment
		}
		if err := tx.Create(appointment).Error; err != nil {
			return err
		}
		result = &CreateAppointmentResult{ID: appointment.ID}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// TODO: prevent fetching old appointments (at least for user)
func (r *Repo) GetAvailability(ctx context.Context, trainerID uint, start, end time.Time) ([]*Appointment, error) {
	existing, err := r.GetAppointments(ctx, trainerID, start, end)
	if err != nil {
		return nil, err
	}

	existingCursor := 0
	availability := make([]*Appointment, 0)
	for opening := range r.GenerateAppointments(start, end) {
		// exclude existing appointments
		hasExisting := false
		for i := existingCursor; i < len(existing); i++ {
			attempt := existing[i]
			if attempt.StartsAt == opening.StartsAt {
				hasExisting = true
				existingCursor++
				break
			}
			// TODO: attempt.EndsAt could theoretically span the next opening's StartsAt (especially if appointment durations change or shift off the 00&30 after cadence)
		}
		if !hasExisting {
			availability = append(availability, opening)
		}
	}

	if len(availability) == 0 {
		return nil, ErrNoAvailability
	}

	return availability, nil
}

func (r *Repo) GetAppointments(ctx context.Context, trainerID uint, start time.Time, end time.Time) ([]*Appointment, error) {
	// TODO: assumption that fetching time out of DB is in UTC, need to verify if time permits
	// TODO: query in UTC
	var appointments []*Appointment
	res := r.db.
		Where(`
			trainer_id = ? AND
			starts_at >= ? AND
			ends_at <= ?`,
			trainerID,
			start.UTC(), //.UTC(),
			end.UTC(),   //.UTC(),
		).
		Find(&appointments).
		Order("starts_at")
	if res.Error != nil {
		return nil, res.Error
	}

	for _, appointment := range appointments {
		appointment.StartsAt = appointment.StartsAt.In(r.location)
		appointment.EndsAt = appointment.EndsAt.In(r.location)
	}

	return appointments, nil
}

func (r *Repo) IsAvailable(ctx context.Context, trainerID uint, start time.Time, end time.Time) (bool, error) {
	var count int64
	// TODO: query in UTC
	res := r.db.Model(&Appointment{}).
		Where(`
			trainer_id = ? AND
			starts_at <= ? AND
			ends_at >= ?`,
			trainerID,
			start.UTC(),
			end.UTC(),
		).
		Count(&count)
	if res.Error != nil {
		return false, res.Error
	}
	return count == 0, nil
}

// Create available appointment between the start and end times.
func (r *Repo) GenerateAppointments(start, end time.Time) <-chan *Appointment {
	return lo.Generator(100, func(yield func(*Appointment)) {
		current := start.In(r.location) // might not be necessary

		for current.Before(end) {
			attempt, err := r.FindOpening(current, end)
			if err != nil {
				return
			}

			current = attempt
			next := current.Add(AppointmentDuration)

			yield(&Appointment{
				StartsAt: current,
				EndsAt:   next,
			})

			current = next
		}
	})
}

func (r *Repo) IsBusinessOpen(t time.Time) bool {
	// UTC Problems:
	// - hours issue: 8AM-5PM PST is 16, 17, 18, 19, 20, 21, 22, 23, 0 UTC (cant compare number between 8 & 17)
	// - weekday issue: 5PM PST would be 00 UTC Saturday (Mon - Fri range won't work)
	// - no daylight savings in UTC
	hour := t.Hour()
	weekday := t.Weekday()
	isBusinessDay := weekday >= time.Monday && weekday <= time.Friday
	isBusinessHour := hour >= BusinessStartHour && hour < BusinessEndHour // TODO: bug where endsat=5pm is invalid
	isOpen := isBusinessDay && isBusinessHour
	return isOpen
}

// Find the next available appointment (will include current if it's a business time).
func (r *Repo) FindOpening(current time.Time, end time.Time) (time.Time, error) {
	current = current.Truncate(time.Minute) // inefficient, move up a level
	end = end.Truncate(time.Minute)
	for current.Before(end) {
		//fmt.Printf("current: %v hour %v weekday %s\n", current, hour, weekday)
		if r.IsBusinessOpen(current) {
			return current, nil
		}

		// TODO: optimizations to skip unnecessary iterations
		// non-business day -> next business day
		// before business hours -> business hours
		// after business hours -> next day
		current = current.Add(AppointmentDuration) // move to next appointment (this is inefficient)
	}

	return time.Time{}, ErrNoAvailability
}
