package repo

import (
	"context"
	"errors"
	"time"

	"github.com/ericbutera/project/internal/db"
	"github.com/ericbutera/project/internal/models"
)

var (
	ErrValidation = errors.New("validation error")
	ErrNotFound   = errors.New("not found")
)

type CreateRecordResult struct {
	ID uint `json:"id"`
}

// Repository pattern for handling business logic.
// Actual database interactions are delegated to the `db` package.
type Repo struct {
	db db.DB
}

func New(db db.DB) (*Repo, error) {
	return &Repo{
		db: db,
	}, nil
}

type StoreReadingsResult struct {
	CreateCount  int
	UpdateCount  int
	InvalidCount int
}

func (r *Repo) StoreReadings(_ context.Context, device string, readings []*models.Reading) (*StoreReadingsResult, error) {
	// TODO: validation should start at the database, then work up through repo and finally into API
	// however for this project I did validation in the API to save time.
	// if err := r.ValidateReadings(device, readings); err != nil {
	// 	return nil, err
	// }
	_, err := r.db.StoreDeviceReadings(device, readings)
	if err != nil {
		return nil, err
	}
	return &StoreReadingsResult{
		// TODO: fill these in
		CreateCount:  0,
		UpdateCount:  0,
		InvalidCount: 0, // TODO: list of invalid readings
	}, nil
}

type GetCountByDeviceResult struct {
	Count int64
}

func (r *Repo) GetCountByDevice(_ context.Context, deviceID string) (*GetCountByDeviceResult, error) {
	res, err := r.db.GetReadingCountByDevice(deviceID)
	if err != nil {
		return nil, err
	}
	return &GetCountByDeviceResult{
		Count: res.Count,
	}, nil
}

type GetLatestReadingByDeviceResult struct {
	LatestReading time.Time
}

func (r *Repo) GetLatestReadingByDevice(_ context.Context, deviceID string) (*GetLatestReadingByDeviceResult, error) {
	// Note: this project only cares about the timestamp, not the full reading data
	// I have left room in the Response for adding in the full `Reading` data if requirements change
	res, err := r.db.GetLatestReadingByDevice(deviceID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) { // TODO: error mapping between layers
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &GetLatestReadingByDeviceResult{
		LatestReading: res.Timestamp,
	}, nil
}

func (r *Repo) GetReadingsByDevice(_ context.Context, deviceID string) ([]models.Reading, error) {
	readings, err := r.db.GetReadingsByDevice(deviceID)
	if err != nil {
		return nil, err
	}
	return readings, nil
}

// TODO: if things get more complex, consider using a validation library
/*
type ValidationError struct {
	Errors []error `json:"errors"`
}

func (e ValidationError) Error() string {
	msgs := make([]string, 0, len(e.Errors))
	for _, err := range e.Errors {
		msgs = append(msgs, err.Error())
	}
	return fmt.Sprintf("validation failed: %s", strings.Join(msgs, "; "))
}

func (r *Repo) ValidateReadings(device string, readings []*models.Reading) error {
	var errs []error
	validate := validator.New()

	if _, err := uuid.Parse(device); err != nil {
		errs = append(errs, errors.New("device ID is not a valid UUID"))
	}

	err := validate.Struct(readings)
	if err != nil {
		var validationErrs validator.ValidationError
		if errors.As(err, &validationErrs) {
			for _, err := range validationErrs {
				errs = append(errs, errors.New(err.Error()))
			}
		}
	}

	if len(errs) > 0 {
		return ValidationError{Errors: errs}
	}

	return nil
}
*/
